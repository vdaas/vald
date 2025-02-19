//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

use futures::{Stream, StreamExt};
use std::sync::Arc;
use tokio::sync::Mutex;
use tonic::{Request, Response, Status, Streaming};
use tokio_stream::wrappers::ReceiverStream;
use tokio::sync::mpsc;

#[macro_export]
macro_rules! stream_type {
    ($t:ty) => {
        std::pin::Pin<Box<dyn tokio_stream::Stream<Item = std::result::Result<$t, tonic::Status>> + Send>>
    };
}

pub async fn bidirectional_stream<Q, R, F, Fut>(
    request_stream: Request<Streaming<Q>>,
    concurrency: usize,
    f: F,
) -> Result<Response<impl Stream<Item = Result<R, Status>>>, Status>
where
    Q: Send + 'static,
    R: Send + 'static,
    F: Fn(Q) -> Fut + Send + Sync + 'static,
    Fut: std::future::Future<Output = Result<R, Status>> + Send + 'static,
{
    let (tx, rx) = mpsc::channel(concurrency);
    let tx = Arc::new(Mutex::new(tx));
    
    let stream = request_stream.into_inner();
    let f = Arc::new(f);

    tokio::spawn(async move {
        let mut handles = Vec::new();
        
        let mut stream = stream;
        while let Some(request) = stream.next().await {
            match request {
                Ok(req) => {
                    let tx = tx.clone();
                    let f = f.clone();
                    
                    let handle = tokio::spawn(async move {
                        let result = f(req).await;
                        let tx = tx.lock().await;
                        let _ = tx.send(result).await;
                    });
                    
                    handles.push(handle);
                    
                    if handles.len() >= concurrency {
                        let done = handles.remove(0);
                        let _ = done.await;
                    }
                }
                Err(e) => {
                    let tx = tx.lock().await;
                    let _ = tx.send(Err(e)).await;
                }
            }
        }

        for handle in handles {
            let _ = handle.await;
        }
    });

    let output_stream = ReceiverStream::new(rx);
    Ok(Response::new(output_stream))
}

#[cfg(test)]
mod tests {
    use super::*;

    use bytes::{Buf, BufMut, Bytes, BytesMut};
    use prost::Message;
    use std::{
        collections::VecDeque,
        marker::PhantomData,
        pin::Pin,
        task::{Context, Poll},
        time::Duration,
    };
    use tokio::time::sleep;
    use tonic::{
        codec::{DecodeBuf, Decoder},
        Status,
    };

    // tonic-mock uses old version of http_body, so we need to implement below ourselves.
    #[derive(Clone)]
    pub struct MockBody {
        data: VecDeque<Bytes>,
    }

    impl MockBody {
        pub fn new(data: Vec<impl Message>) -> Self {
            let mut queue: VecDeque<Bytes> = VecDeque::with_capacity(16);
            for msg in data {
                let buf = Self::encode(msg);
                queue.push_back(buf);
            }

            MockBody { data: queue }
        }

        pub fn is_empty(&self) -> bool {
            self.data.is_empty()
        }

        // see: https://github.com/hyperium/tonic/blob/1b03ece2a81cb7e8b1922b3c3c1f496bd402d76c/tonic/src/codec/encode.rs#L52
        fn encode(msg: impl Message) -> Bytes {
            let mut buf = BytesMut::with_capacity(256);

            buf.reserve(5);
            unsafe {
                buf.advance_mut(5);
            }
            msg.encode(&mut buf).unwrap();
            {
                let len = buf.len() - 5;
                let mut buf = &mut buf[..5];
                buf.put_u8(0); // byte must be 0, reserve doesn't auto-zero
                buf.put_u32(len as u32);
            }
            buf.freeze()
        }
    }

    impl http_body::Body for MockBody {
        type Data = Bytes;
        type Error = Status;

        fn poll_frame(
            mut self: Pin<&mut Self>,
            _cx: &mut Context<'_>,
        ) -> Poll<Option<Result<http_body::Frame<Self::Data>, Self::Error>>> {
            if !self.is_empty() {
                let data = self.data.pop_front().unwrap();
                Poll::Ready(Some(Ok(http_body::Frame::data(data))))
            } else {
                Poll::Ready(None)
            }
        }
    
        fn is_end_stream(&self) -> bool {
            self.is_empty()
        }
    
        fn size_hint(&self) -> http_body::SizeHint {
            let mut hint = http_body::SizeHint::new();
            let remaining = self.data.iter().map(|b| b.len()).sum::<usize>();
            hint.set_exact(remaining as u64);
            hint
        }
    }

    #[derive(Debug, Clone, Default)]
    pub struct ProstDecoder<U>(PhantomData<U>);

    impl<U> ProstDecoder<U> {
        pub fn new() -> Self {
            Self(PhantomData)
        }
    }

    impl<U: Message + Default> Decoder for ProstDecoder<U> {
        type Item = U;
        type Error = Status;

        fn decode(&mut self, buf: &mut DecodeBuf<'_>) -> Result<Option<Self::Item>, Self::Error> {
            let item = Message::decode(buf.chunk())
                .map(Option::Some)
                .map_err(|e| Status::internal(e.to_string()))?;

            buf.advance(buf.chunk().len());
            Ok(item)
        }
    }

    #[tokio::test]
    async fn test_bidirectional_stream() {
        // テストデータの準備
        let decoder: ProstDecoder<i32> = ProstDecoder::new();
        let messages = vec![1, 2, 3, 4, 5];
        let body = MockBody::new(messages);
        let streaming = Streaming::new_request(decoder, body, None, None);
        let request = Request::new(streaming);
        
        // テスト用の処理関数
        let process_fn = |n: i32| async move {
            sleep(Duration::from_millis(10)).await;
            Ok(n * 2)
        };

        // bidirectional_streamの実行
        let response = bidirectional_stream(request, 3, process_fn).await.unwrap();
        let mut stream = response.into_inner();

        let mut results = Vec::new();
        while let Some(result) = stream.next().await {
            match result {
                Ok(n) => results.push(n),
                Err(_) => break,
            }
        }

        assert_eq!(results, vec![2, 4, 6, 8, 10]);
    }

    #[tokio::test]
    async fn test_bidirectional_stream_with_error() {
        let decoder: ProstDecoder<i32> = ProstDecoder::new();
        let messages = vec![1, 2, 3, 4, 5];
        let body = MockBody::new(messages);
        let streaming = Streaming::new_request(decoder, body, None, None);
        let request = Request::new(streaming);

        let process_fn = |n: i32| async move {
            if n == 3 {
                Err(Status::internal("test error"))
            } else {
                Ok(n * 2)
            }
        };

        let response = bidirectional_stream(request, 2, process_fn).await.unwrap();
        let mut stream = response.into_inner();

        let mut results = Vec::new();
        let mut errors = Vec::new();
        while let Some(result) = stream.next().await {
            match result {
                Ok(n) => results.push(n),
                Err(e) => errors.push(e.message().to_string()),
            }
        }

        assert_eq!(results, vec![2, 4, 8, 10]);
        assert_eq!(errors, vec!["test error"]);
    }
}
