# Use cases

This document will introduce you to the example of what Vald can do.
Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine, which uses [NGT](https://github.com/yahoojapan/NGT) as the core engine of Vald, and Vald manages to integrate with Kubernetes.

In general, you cannot search your unstructured data using the inverted index, like images and videos.
By applying the model like BERT or VGG, you can convert your unstructured data into vectors.
After converting them into vectors, you can insert it to Vald and process them in Vald.

Here are some general use cases of Vald or vector search engine.

## Image and video processing

  You can use Vald as the image/video processing engine to search the similar image/video or analysis the image/video for your use case.

  Vald is capable to process a huge number of images at the same time so it case fit with your use case.

  Here are some examples of what you can do with images and videos in Vald.

  1. Search by image
  1. Face recognition
  1. Product recommendation base on images
  1. Image/Video analysis
  1. Image/Video deduplication

## Audio processing

  Audio processing is important for personal assistant implementation.

  Vald can act as a brain of the personal assistant function, conversation interpreter and the natural language generation.

  Here are some examples of what you can process in Vald.

  1. Personal assistant
  1. Speech recognition
  1. Natural language understanding and generation

## Text processing

  Since Vald supports some texting vectorizing like `BERT`, you can process your text in Vald.

  Here are some examples of the use case of text processing in Vald.

  1. Search by text
  1. Product recommendation based on text
  1. Grammar checker
  1. Real-time translator

## Data analysis

  Vald can process the vector data, you can analyze every data you can vectorize.

  Here are some examples of the use case of data analysis.

  1. AI malware detection
  1. Price optimization
  1. Social analysis

## Advanced usecase

  Besides the general use case of Vald or vector search engine, Vald supports a user-defined filter that the user can customize the filter to filter the specific result.

  For example when the user chose a man's t-shirt and the recommended product is going to be searched in Vald.

  Without the filtering functionality, the women's t-shirt may be searched in Vald and displayed because women's t-shirt is similar to the men's t-shirt and it is very hard to differentiate the image of men's and women's t-shirt.

  By implementing the custom filter, you can filter only the man's t-shirt based on your criteria and needs.
