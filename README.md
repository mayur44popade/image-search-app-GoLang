# Image Search Application

This is a very simple image search application. It takes a word as input and returns top 10 best matching images for that word.
For example, if you enter 'dog' application will display images of dog. I have used a set of 1000 images for this application. The list is provided in imagesAll.txt
You can use images.txt for testing purposes. I have used Clarifai API to get keywords associated with these images and then I stored these words in application memory.

## Getting Started

To run this application on your local system, first clone this repository and run following command:
'go run forms.go'
If you are testing on images.txt, the application will be up and running in a minute.
If you run it using imagesAll.txt, it would take around 10-12 minutes to collect data from Clarifai and start running.

### Prerequisites

You need to install Go to run this application.

## Deployment

You can also deply this on server if you want

## Built With

* JetBrains GoLand 2018.2.4

## Authors

* **Mayur Popade** - *Initial work* - [mayur44popade](https://github.com/mayur44popade)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
