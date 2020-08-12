# Use cases

  This document will introduce you to the example of what Vald can do.
  Vald is a highly scalable distributed fast approximate nearest neighbor dense vector search engine, which uses [NGT](https://github.com/yahoojapan/NGT) as the core engine of Vald, and Vald manages to integrate with Kubernetes.

  In general, you cannot search your unstructured data using the inverted index, like images and videos.
  By applying the model like BERT or VGG, you can convert your unstructured data into vectors.
  After converting them into vectors, you can insert it to Vald and process them in Vald.

  Here are some general use cases of Vald or vector search engine.

## Image and video processing

  <img src="../../assets/docs/usecase_image.png" />

  You can use Vald as the image/video processing engine to search the similar image/video or analyze the image/video for your use case.

  Vald is capable to process a huge number of images at the same time so it case fit with your use case.

  Here are some examples of what you can do with images and videos using Vald.

- Search by image
- Face recognition
- Product recommendation based on images
- Image/Video analysis
- Image/Video deduplication

## Audio processing

  <img src="../../assets/docs/usecase_audio.png" />

  Audio processing is important for personal assistant implementation.

  Vald can act as a brain of the personal assistant function, conversation interpreter and the natural language generation.

  Here are some examples of what you can process using Vald.

- Personal assistant
- Speech recognition
- Natural language understanding and generation

## Text processing

  <img src="../../assets/docs/usecase_text.png" />

  Using a text vectorizing model like BERT, you can process your text data in Vald.

  Here are some examples of the use case of text processing using Vald.

- Search by text
- Product recommendation based on text
- Grammar checker
- Real-time translator

## Data analysis

  <img src="../../assets/docs/usecase_data.png" />

  Vald can process the vector data, you can analyze every data you can vectorize.

  Here are some examples of the use case of data analysis.

- AI malware detection

  To detect the malware using Vald, you need to vectorize the malware binary file and insert into Vald first.
  You can analyze your binary by performing a search to find a similar binary in Vald.
  If your binary is similar to the malware binary, you can trigger the alert for users.

- Price optimization

   By applying the price optimization technique using Vald, you can find the most optimized price for your business.
  You can apply models like GLMs to achieve it and use Vald as a machine learning engine for your business.

- Social analysis

  To analyze the social relationship of users, you can suggest them their related friends, page recommendation, or other use cases.
  You can apply different models to analyze the social data, and use Vald as a recommendation engine for your business.

## Advanced use cases

  Besides the general use case of Vald or vector search engine, Vald supports a user-defined filter that the user can customize the filter to filter the specific result.

  For example when the user chose a man's t-shirt and the recommended product is going to be searched in Vald.

  Without the filtering functionality, the women's t-shirt may be searched in Vald and displayed because women's t-shirt is similar to the men's t-shirt and it is very hard to differentiate the image of men's and women's t-shirt.

  By implementing the custom filter, you can filter only the man's t-shirt based on your criteria and needs.
