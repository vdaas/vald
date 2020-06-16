# Use cases

This document will introduce you the example of what Vald can do.
Vald is the approximate nearest neighbor dense vector search engine. Vald supports the general usecase of similarity search engine, like faiss.

Here is the example of the use cases of Vald.

1. Image and video processing

    You can use Vald as the image/video processing engine to search the similar image/video or analysis the image/video for your use case.

    Vald is capable to process huge number of images at the same time so it should fit with your use case.

    Here are some example of what you can do with images and video in Vald.

    1. Search by image
    1. Face recognition
    1. Product recommandation base on images
    1. Image/Video analysis
    1. Image/Video deduplication

    In addition Vald support user-defined filter that user can customize the filter to filter the specific result.

    For example when user chose a man's t-shirt and the similar product is going to be searched in Vald.

    Without the filtering functionality the women's t-shirt may be searched in Vald and displayed because women's t-shirt is similar to the men's t-shirt and it is very hard to differentiate the image of men's and women's t-shirt.

    By implementing the custom filter you can filter only the man's t-shirt base on your criteria.

1. Audio processing

    Audio processing is important for the personal assistant implementation.

    Vald can act as a brain of the personal assistant function, conversation interpretor and the natural language generation.

    Here are some example of what you can process in Vald.

    1. Personal assistant
    1. Speech recognition
    1. Natural language understanding and generation

1. Text processing

    Since Vald support some texting vectorizing engine like `BERT`, you can process your text in Vald.

    Here are some example of the use case of the text processing in Vald.

    1. Search by text
    1. Product recommandation base on text
    1. Grammar checker
    1. Real-time translator

1. Data Analysis

    Vald can process the vector data, you can analysis every data you can vectorize.

    Here are some example of the use case of data analysis.

    1. AI malware detection
    1. Price optimization
    1. Social analysis

