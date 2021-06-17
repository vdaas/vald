//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package uploader

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/arn"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/pool"
	"github.com/vdaas/vald/internal/db/storage/blob/v2/s3/ua"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
)

// MaxUploadParts is the maximum allowed number of parts in a multi-part upload
// on Amazon S3.
const MaxUploadParts int32 = 10000

// MinUploadPartSize is the minimum allowed part size when uploading a part to
// Amazon S3.
const MinUploadPartSize int64 = 1024 * 1024 * 5

// DefaultUploadPartSize is the default part size to buffer chunks of a
// payload into.
const DefaultUploadPartSize = MinUploadPartSize

// DefaultUploadConcurrency is the default number of goroutines to spin up when
// using Upload().
const DefaultUploadConcurrency = 5

// A MultiUploadFailure wraps a failed S3 multipart upload. An error returned
// will satisfy this interface when a multi part upload failed to upload all
// chucks to S3. In the case of a failure the UploadID is needed to operate on
// the chunks, if any, which were uploaded.
//
// Example:
//
//	u := manager.NewClient(client)
//	output, err := u.upload(context.Background(), input)
//	if err != nil {
//		var multierr manager.MultiUploadFailure
//		if errors.As(err, &multierr) {
//			fmt.Printf("upload failure UploadID=%s, %s\n", multierr.UploadID(), multierr.Error())
//		} else {
//			fmt.Printf("upload failure, %s\n", err.Error())
//		}
//	}
//
type MultiUploadFailure interface {
	error

	// UploadID returns the upload id for the S3 multipart upload that failed.
	UploadID() string
}

// A multiUploadError wraps the upload ID of a failed s3 multipart upload.
// Composed of BaseError for code, message, and original error
//
// Should be used for an error that occurred failing a S3 multipart upload,
// and a upload ID is available. If an uploadID is not available a more relevant
type multiUploadError struct {
	err error

	// ID for multipart upload which failed.
	uploadID string
}

// batchItemError returns the string representation of the error.
//
// See apierr.BaseError ErrorWithExtra for output format
//
// Satisfies the error interface.
func (m *multiUploadError) Error() string {
	var extra string
	if m.err != nil {
		extra = fmt.Sprintf(", cause: %s", m.err.Error())
	}
	return fmt.Sprintf("upload multipart failed, upload id: %s%s", m.uploadID, extra)
}

// Unwrap returns the underlying error that cause the upload failure
func (m *multiUploadError) Unwrap() error {
	return m.err
}

// UploadID returns the id of the S3 upload which failed.
func (m *multiUploadError) UploadID() string {
	return m.uploadID
}

// UploadOutput represents a response from the Upload() call.
type UploadOutput struct {
	// The URL where the object was uploaded to.
	Location string

	// The version of the object that was uploaded. Will only be populated if
	// the S3 Bucket is versioned. If the bucket is not versioned this field
	// will not be set.
	VersionID *string

	// The ID for a multipart upload to S3. In the case of an error the error
	// can be cast to the MultiUploadFailure interface to extract the upload ID.
	UploadID string
}

// WithUploadRequestOptions appends to the Client's API client options.
func WithUploadRequestOptions(opts ...func(*s3.Options)) func(*Client) {
	return func(u *Client) {
		u.ClientOptions = append(u.ClientOptions, opts...)
	}
}

// The Client structure that calls Upload(). It is safe to call Upload()
// on this structure for multiple objects and across concurrent goroutines.
// Mutating the Client's properties is not safe to be done concurrently.
type Client struct {
	// The buffer size (in bytes) to use when buffering data into chunks and
	// sending them as parts to S3. The minimum allowed part size is 5MB, and
	// if this value is set to zero, the DefaultUploadPartSize value will be used.
	PartSize int64

	// The number of goroutines to spin up in parallel per call to Upload when
	// sending parts. If this is set to zero, the DefaultUploadConcurrency value
	// will be used.
	//
	// The concurrency pool is not shared between calls to Upload.
	Concurrency int

	// Setting this value to true will cause the SDK to avoid calling
	// AbortMultipartUpload on a failure, leaving all successfully uploaded
	// parts on S3 for manual recovery.
	//
	// Note that storing parts of an incomplete multipart upload counts towards
	// space usage on S3 and will add additional costs if not cleaned up.
	LeavePartsOnError bool

	// MaxUploadParts is the max number of parts which will be uploaded to S3.
	// Will be used to calculate the partsize of the object to be uploaded.
	// E.g: 5GB file, with MaxUploadParts set to 100, will upload the file
	// as 100, 50MB parts. With a limited of s3.MaxUploadParts (10,000 parts).
	//
	// MaxUploadParts must not be used to limit the total number of bytes uploaded.
	// Use a type like to io.LimitReader (https://golang.org/pkg/io/#LimitedReader)
	// instead. An io.LimitReader is helpful when uploading an unbounded reader
	// to S3, and you know its maximum size. Otherwise the reader's io.EOF returned
	// error must be used to signal end of stream.
	//
	// Defaults to package const's MaxUploadParts value.
	MaxUploadParts int32

	// The client to use when uploading to S3.
	S3 manager.UploadAPIClient

	// List of request options that will be passed down to individual API
	// operation requests made by the uploader.
	ClientOptions []func(*s3.Options)

	// Defines the buffer strategy used when uploading a part
	BufferProvider manager.ReadSeekerWriteToProvider

	// partPool allows for the re-usage of streaming payload part buffers between upload calls
	partPool pool.ByteSlicePool
}

func New(client manager.UploadAPIClient, options ...func(*Client)) *Client {
	u := &Client{
		S3:                client,
		PartSize:          DefaultUploadPartSize,
		Concurrency:       DefaultUploadConcurrency,
		LeavePartsOnError: false,
		MaxUploadParts:    MaxUploadParts,
		BufferProvider:    nil,
	}

	for _, option := range options {
		option(u)
	}

	u.partPool = pool.NewByteSlicePool(u.PartSize)

	return u
}

// Upload uploads an object to S3, intelligently buffering large
// files into smaller chunks and sending them in parallel across multiple
// goroutines. You can configure the buffer size and concurrency through the
// Client parameters.
//
// Additional functional options can be provided to configure the individual
// upload. These options are copies of the Client instance Upload is called from.
// Modifying the options will not impact the original Client instance.
//
// Use the WithClientRequestOptions helper function to pass in request
// options that will be applied to all API operations made with this uploader.
//
// It is safe to call this method concurrently across goroutines.
func (u Client) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*Client)) (*UploadOutput, error) {
	i := uploader{in: input, cfg: u, ctx: ctx}

	// Copy ClientOptions
	clientOptions := make([]func(*s3.Options), 0, len(i.cfg.ClientOptions)+1)
	clientOptions = append(clientOptions, func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions, middleware.AddSDKAgentKey(middleware.FeatureMetadata, ua.Key))
	})
	clientOptions = append(clientOptions, i.cfg.ClientOptions...)
	i.cfg.ClientOptions = clientOptions

	for _, opt := range opts {
		opt(&i.cfg)
	}

	return i.upload()
}

// internal structure to manage an upload to S3.
type uploader struct {
	ctx context.Context
	cfg Client

	in *s3.PutObjectInput

	readerPos int64 // current reader position
	totalSize int64 // set to -1 if the size is not known
}

// internal logic for deciding whether to upload a single part or use a
// multipart upload.
func (u *uploader) upload() (*UploadOutput, error) {
	if err := arn.ValidateSupportedARNType(aws.ToString(u.in.Bucket)); err != nil {
		return nil, fmt.Errorf("unable to initialize upload: %w", err)
	}

	if u.cfg.Concurrency == 0 {
		u.cfg.Concurrency = DefaultUploadConcurrency
	}
	if u.cfg.PartSize == 0 {
		u.cfg.PartSize = DefaultUploadPartSize
	}
	if u.cfg.MaxUploadParts == 0 {
		u.cfg.MaxUploadParts = MaxUploadParts
	}

	u.totalSize = -1

	switch r := u.in.Body.(type) {
	case io.Seeker:
		curOffset, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize upload: %w", err)
		}

		endOffset, err := r.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize upload: %w", err)
		}

		_, err = r.Seek(curOffset, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize upload: %w", err)
		}

		u.totalSize = endOffset - curOffset
		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		if u.totalSize/u.cfg.PartSize >= int64(u.cfg.MaxUploadParts) {
			// Add one to the part size to account for remainders
			// during the size calculation. e.g odd number of bytes.
			u.cfg.PartSize = (u.totalSize / int64(u.cfg.MaxUploadParts)) + 1
		}
	}
	// If PartSize was changed or partPool was never setup then we need to allocated a new pool
	// so that we return []byte slices of the correct size
	poolCap := u.cfg.Concurrency + 1
	if u.cfg.partPool == nil || u.cfg.partPool.SliceSize() != u.cfg.PartSize {
		u.cfg.partPool = pool.NewByteSlicePool(u.cfg.PartSize)
	} else {
		u.cfg.partPool = &pool.ReturnCapacityPoolCloser{ByteSlicePool: u.cfg.partPool}
	}
	u.cfg.partPool.ModifyCapacity(poolCap)
	defer u.cfg.partPool.Close()

	if u.cfg.PartSize < MinUploadPartSize {
		return nil, fmt.Errorf("part size must be at least %d bytes", MinUploadPartSize)
	}

	// Do one read to determine if we have more than one part
	reader, _, cleanup, err := u.nextReader()
	if err == io.EOF { // single part
		defer cleanup()
		params := u.in
		params.Body = reader
		var locationRecorder recordLocationClient
		out, err := u.cfg.S3.PutObject(u.ctx, params, append(u.cfg.ClientOptions, locationRecorder.WrapClient())...)
		if err != nil {
			return nil, err
		}

		return &UploadOutput{
			Location:  locationRecorder.location,
			VersionID: out.VersionId,
		}, nil
	} else if err != nil {
		cleanup()
		return nil, fmt.Errorf("read upload data failed: %w", err)
	}
	mu := multiuploader{uploader: u}
	return mu.upload(reader, cleanup)
}

func (u *uploader) nextReader() (io.ReadSeeker, int, func(), error) {
	switch r := u.in.Body.(type) {
	case readerAtSeeker:
		var err error

		n := u.cfg.PartSize
		if u.totalSize >= 0 {
			bytesLeft := u.totalSize - u.readerPos

			if bytesLeft <= u.cfg.PartSize {
				err = io.EOF
				n = bytesLeft
			}
		}

		var (
			reader  io.ReadSeeker
			cleanup func()
		)

		reader = io.NewSectionReader(r, u.readerPos, n)
		if u.cfg.BufferProvider != nil {
			reader, cleanup = u.cfg.BufferProvider.GetWriteTo(reader)
		} else {
			cleanup = func() {}
		}

		u.readerPos += n

		return reader, int(n), cleanup, err

	default:
		part, err := u.cfg.partPool.Get(u.ctx)
		if err != nil {
			return nil, 0, func() {}, err
		}
		var offset int
		var n int
		b := *part
		for offset < len(b) && err == nil {
			var n int
			n, err = r.Read(b[offset:])
			offset += n
		}
		u.readerPos += int64(n)

		cleanup := func() {
			u.cfg.partPool.Put(part)
		}

		return bytes.NewReader((*part)[0:n]), n, cleanup, err
	}
}

type httpClient interface {
	Do(r *http.Request) (*http.Response, error)
}

type recordLocationClient struct {
	httpClient
	location string
}

func (c *recordLocationClient) WrapClient() func(o *s3.Options) {
	return func(o *s3.Options) {
		c.httpClient = o.HTTPClient
		o.HTTPClient = c
	}
}

func (c *recordLocationClient) Do(r *http.Request) (resp *http.Response, err error) {
	resp, err = c.httpClient.Do(r)
	if err != nil {
		return resp, err
	}

	if resp.Request != nil && resp.Request.URL != nil {
		url := *resp.Request.URL
		url.RawQuery = ""
		c.location = url.String()
	}

	return resp, err
}

// internal structure to manage a specific multipart upload to S3.
type multiuploader struct {
	*uploader
	wg       sync.WaitGroup
	m        sync.Mutex
	err      error
	uploadID string
	parts    completedParts
}

// keeps track of a single chunk of data being sent to S3.
type chunk struct {
	buf     io.ReadSeeker
	num     int32
	cleanup func()
}

// completedParts is a wrapper to make parts sortable by their part number,
// since S3 required this list to be sent in sorted order.
type completedParts []types.CompletedPart

func (a completedParts) Len() int           { return len(a) }
func (a completedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a completedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

// upload will perform a multipart upload using the firstBuf buffer containing
// the first chunk of data.
func (u *multiuploader) upload(firstBuf io.ReadSeeker, cleanup func()) (*UploadOutput, error) {
	// Create the multipart
	var locationRecorder recordLocationClient
	resp, err := u.cfg.S3.CreateMultipartUpload(u.ctx, u.in, append(u.cfg.ClientOptions, locationRecorder.WrapClient())...)
	if err != nil {
		cleanup()
		return nil, err
	}
	u.uploadID = *resp.UploadId

	// Create the workers
	ch := make(chan chunk, u.cfg.Concurrency)
	for i := 0; i < u.cfg.Concurrency; i++ {
		u.wg.Add(1)
		go func(ch chan chunk) {
			defer u.wg.Done()
			for {
				data, ok := <-ch

				if !ok {
					break
				}

				if u.geterr() == nil {
					err := func(c chunk) error {
						resp, err := u.cfg.S3.UploadPart(u.ctx,
							&s3.UploadPartInput{
								Bucket:               u.in.Bucket,
								Key:                  u.in.Key,
								Body:                 c.buf,
								UploadId:             &u.uploadID,
								SSECustomerAlgorithm: u.in.SSECustomerAlgorithm,
								SSECustomerKey:       u.in.SSECustomerKey,
								PartNumber:           c.num,
							}, u.cfg.ClientOptions...)
						if err != nil {
							return err
						}

						n := c.num
						completed := types.CompletedPart{ETag: resp.ETag, PartNumber: n}

						u.m.Lock()
						u.parts = append(u.parts, completed)
						u.m.Unlock()

						return nil
					}(data)
					if err != nil {
						u.seterr(err)
					}
				}

				data.cleanup()
			}
		}(ch)
	}

	// Send part 1 to the workers
	var num int32 = 1
	ch <- chunk{buf: firstBuf, num: num, cleanup: cleanup}

	// Read and queue the rest of the parts
	for u.geterr() == nil && err == nil {
		var (
			reader       io.ReadSeeker
			nextChunkLen int
		)

		reader, nextChunkLen, cleanup, err = u.nextReader()

		if err != nil && err != io.EOF {
			cleanup()
			if err != nil {
				u.seterr(fmt.Errorf("read multipart upload data failed, %w", err))
			}
			break
		}

		if nextChunkLen == 0 {
			cleanup()
			break
		}

		num++
		// This upload exceeded maximum number of supported parts, error now.
		if num > u.cfg.MaxUploadParts || num > MaxUploadParts {
			var msg string
			if num > u.cfg.MaxUploadParts {
				msg = fmt.Sprintf("exceeded total allowed configured MaxUploadParts (%d). Adjust PartSize to fit in this limit",
					u.cfg.MaxUploadParts)
			} else {
				msg = fmt.Sprintf("exceeded total allowed S3 limit MaxUploadParts (%d). Adjust PartSize to fit in this limit",
					MaxUploadParts)
			}
			cleanup()
			if err != nil {
				u.seterr(fmt.Errorf(msg))
			}
			break
		}

		num++

		ch <- chunk{buf: reader, num: num, cleanup: cleanup}
	}

	close(ch)
	u.wg.Wait()

	// Parts must be sorted in PartNumber order.
	sort.Sort(u.parts)

	params := &s3.CompleteMultipartUploadInput{
		Bucket:          u.in.Bucket,
		Key:             u.in.Key,
		UploadId:        &u.uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{Parts: u.parts},
	}
	complete, err := u.cfg.S3.CompleteMultipartUpload(u.ctx, params, u.cfg.ClientOptions...)
	if err != nil {
		if !u.cfg.LeavePartsOnError {
			_, err := u.cfg.S3.AbortMultipartUpload(u.ctx, &s3.AbortMultipartUploadInput{
				Bucket:   u.in.Bucket,
				Key:      u.in.Key,
				UploadId: &u.uploadID,
			}, u.cfg.ClientOptions...)
			if err != nil {
				log.Warn(err)
			}
		}
		return nil, &multiUploadError{
			err:      err,
			uploadID: u.uploadID,
		}
	}

	return &UploadOutput{
		Location:  locationRecorder.location,
		VersionID: complete.VersionId,
		UploadID:  u.uploadID,
	}, nil
}

// geterr is a thread-safe getter for the error object
func (u *multiuploader) geterr() error {
	u.m.Lock()
	defer u.m.Unlock()

	return u.err
}

// seterr is a thread-safe setter for the error object
func (u *multiuploader) seterr(e error) {
	u.m.Lock()
	defer u.m.Unlock()

	u.err = e
}

type readerAtSeeker interface {
	io.ReaderAt
	io.ReadSeeker
}

func seekerLen(s io.Seeker) (int64, error) {
	curOffset, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	endOffset, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(curOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}
