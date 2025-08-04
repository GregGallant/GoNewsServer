
# GoNewsServer
![Image](https://github.com/user-attachments/assets/ac1c86e8-e3c5-411f-87a0-b853e79cf929)

## Go News Server
This is only the Go code for an automated Go powered news API parser for Webz.io data.
Webz.io is a news content API, however it did not have any support for the Go language.  

This is a Go news server I wrote to handle any and all webz.io content along with a chi http server to route.  This can be used in a docker container or run bare metal.  Originally I ran the binary as a System-D process but now I just use Docker.  This code is not the most up to date code as it does not have multi-news interfaces, it is just the base representation of the code used in production.

The news object I'm calling in article.go is being used as a receiver for a few orgazational methods I've been frequently updating: a quick algorithm to remove duplicate articles, image url handlers to clean UTF-8 links which can vary greatly depending on the CMS of news aggregates and external websites, etc.  This is how I'm running it on my server, the version here will have to be updated depending on your file system or docker container.

The news file checks for a date file and writes the latest news based upon an filterable query.  You can use a customizable query, write your own AI parser or just use a basic content parser to add to your news server.

The most recent addition was an AI starting point using Prediction Guard's(https://predictionguard.com/) embedded models for LLMs, with an update to the EmbedInputType interface that's currently used in their latest version.  This was only used for testing and understanding how to use vectorized chunks, cosine similarity and a launch pad to use Google's Vertex AI.

You can view it in action here using a React Router frontend: https://www.gallantone.com/news -and- https://www.cronvega.com/news 
