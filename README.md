# GoNewsServer
Automated Go powered news API parser for Webz.io data

## Go News Server
Webz.io is a news content API, however it did not have any support for the Go language.  

This is a Go news server I wrote to handle any and all webz.io content along with a chi http server to route.  This can be used in a docker container or run bare metal.  Originally I ran the binary as a System-D process but now I just use Docker. 

The news file checks for a date file and writes the latest news based upon an filterable query.  You can use a customizable query, write your own AI parser or just use a basic content parser to add to your news server.

You can view it in action here using a Svelte frontend: https://www.gallantone.com/news 
