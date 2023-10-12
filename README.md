<br/>
<p align="center">
  <a href="https://github.com/Xumeiquer/tmd">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Telegram Media Downloader</h3>

  <p align="center">
    Download Telegram media like a pro!
    <br/>
    <br/>
    <a href="https://github.com/Xumeiquer/tmd"><strong>Explore the docs »</strong></a>
    <br/>
    <br/>
    <a href="https://github.com/Xumeiquer/tmd/issues">Report Bug</a>
    .
    <a href="https://github.com/Xumeiquer/tmd/issues">Request Feature</a>
  </p>
</p>

![Downloads](https://img.shields.io/github/downloads/Xumeiquer/tmd/total) ![Issues](https://img.shields.io/github/issues/Xumeiquer/tmd) ![License](https://img.shields.io/github/license/Xumeiquer/tmd) 

## Table Of Contents

* [About the Project](#about-the-project)
* [Built With](#built-with)
* [Getting Started](#getting-started)
  * [Prerequisites](#prerequisites)
* [Usage](#usage)
* [Roadmap](#roadmap)
* [Contributing](#contributing)
* [License](#license)
* [Authors](#authors)
* [Acknowledgements](#acknowledgements)

## About The Project

![Screen Shot](images/screenshot.png)

Telegram Media Downloader can be considered a lightweight Telegram client which allows you to download media files from the Telegram cloud.

## Built With

This project is build using [zelenin](https://github.com/zelenin/go-tdlib) wrapper for the [Telegram Database Library](https://github.com/tdlib/td).

## Getting Started

Telegram Media Downloader is a command line application that runs inside a Docker container. The main reason is to ease the build process.

### Prerequisites

You will need to [register an application](https://my.telegram.org/apps) in Telegram to obtain an `API_ID` and `API_HASH`, Telegram Media Downloader will need the API information to connect to Telegram on your behalf.
On the other hand, you need to have Docker installed to be able to execute Telegram Media Downloader.

## Usage

```shell
$ docker run --rm -it ghcr.io/xumeiquer/tmd --help
Download Telegram media from Users, Chats, Channels, or Forums.

Telegram Media Downloader allow users to download media content from Telegram cloud
without manually interacting with the Telegram client. Telegram Media Downloader is
or acts as a client so it has to be enrolled as a Telegram client.

Usage:
  tmd [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  download    Download media from the Telegram cloud
  help        Help about any command
  list        List information about conversations.

Flags:
  -h, --help               help for tmd
      --log-level string   set log level (info, error, warn, debug)
      --log-to string      where to log (default "stdout")
      --log-type string    log as text or JSON (default "json")
  -v, --version            version for tmd

Use "tmd [command] --help" for more information about a command.
```

There are two main commands `list` and `download`.

### List commands

The `list` commands are intended to search among your chats, channels, forums, etc so you can select what to download.

If you want to list all your chats, execute:

`docker run --rm -it -e TMD_TDLIB_API_ID=$TMD_TDLIB_API_ID -e TMD_TDLIB_API_HASH=$TMD_TDLIB_API_HASH -v ${PWD}/tdlib:/.tdlib -v ${PWD}/media:/media ghcr.io/xumeiquer/tmd:latest list`

## Roadmap

See the [open issues](https://github.com/Xumeiquer/tmd/issues) for a list of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.
* If you have suggestions for adding or removing projects, feel free to [open an issue](https://github.com/Xumeiquer/tmd/issues/new) to discuss it, or directly create a pull request after you edit the *README.md* file with necessary changes.
* Please make sure you check your spelling and grammar.
* Create individual PR for each suggestion.
* Please also read through the [Code Of Conduct](https://github.com/Xumeiquer/tmd/blob/main/CODE_OF_CONDUCT.md) before posting your first idea as well.

### Creating A Pull Request

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See [LICENSE](https://github.com/Xumeiquer/tmd/blob/main/LICENSE.md) for more information.

## Authors

* **Xumeiquer** - ** - [Xumeiquer](https://github.com/Xumeiquer) - **

## Acknowledgements

* []()
* []()
* []()
