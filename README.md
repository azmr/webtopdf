webtopdf
========

Web app to create pdf from images found from provided URL

`go run imgCollate.go` runs small server on localhost:8080, on which there is a form with a 'URL' textbox and a 'Format' dropdown option. When you click 'Get PDF', the app downloads all the (currently just jpg) images from that URL to a temporary directory, makes and serves a pdf of those images in the format specified, and deletes the temporary folder.
    
