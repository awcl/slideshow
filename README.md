# Slideshow

A Go-based tool for creating an image slideshow.

## Overview

The Slideshow application displays images from a specified directory in a web browser as a dynamic slideshow. It is built with Go and provides an easy way to present images in a rotating format.

## Getting Started

0. **Ensure Go is Installed**

   https://go.dev/

1. **Build and Run the Application**

   To build and run the application, simply double-click the `buildrun.bat` file. This script will compile the `slideshow.go` file into an executable named `run.exe` and start it. After running the script, the terminal will display the URL where the slideshow can be accessed (e.g., `http://localhost:3000` or `http://<IP Address>:3000`).

## Usage

1. **Prepare the Image Directory**

   Place the images you want to display in the `/photos` directory. The application supports common image formats such as JPG, JPEG, PNG, and GIF.

2. **Access the Slideshow**

   Open a web browser and navigate to the URL provided by the terminal to view the slideshow. The default URL will be `http://localhost:3000`, but it can also be accessed via your local network IP address.

## Technical Details

- **Supported Image Formats:** JPG, JPEG, PNG, GIF
- **Defaults to JavaScript in HTML:** Optionally edit the `buildrun.bat` to use `nojs_slideshow.go` if display devices aren't JavaScript enabled.
- **Slideshow Interval:** Configurable in the `slideshow.go` or `nojs_slideshow.go` file using the `DELAY_IN_SECONDS` constant (default: 4 seconds).

## Troubleshooting

If the application does not build or run as expected, ensure that:
- Your Go environment is correctly set up.
- The `/photos` directory exists and contains valid image files.
