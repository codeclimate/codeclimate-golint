# Code Climate Golint Engine

`codeclimate-golint` is a Code Climate engine that wraps [Golint](https://github.com/golang/lint). You can run it on your command line using the Code Climate CLI, or on our hosted analysis platform.

Golint is a linter for Go code. Where as [Gofmt](https://www.github.com/codeclimate/codeclimate-gofmt) automatically reformats code, Golint suggests style issues that may need to be addressed.

### Installation

1. If you haven't already, [install the Code Climate CLI](https://github.com/codeclimate/codeclimate).
2. Run `codeclimate engines:enable golint`. This command both installs the engine and enables it in your `.codeclimate.yml` file.
3. You're ready to analyze! Browse into your project's folder and run `codeclimate analyze`.

### Configuration

Like the `golint` binary, you can configure the minimum confidence threshold of
this engine: issues reported by `golint` must have a confidence value higher than
the threshold in order to be reported. The default value is `0.8`, the same as
`golint`: you can configure your own threshold in your `.codeclimate.yml`:

```yaml
engines:
  golint:
    enabled: true
    config:
      min_confidence: 0.1
```

### Building

```console
docker build -t codeclimate/codeclimate-golint .
```

### Updating

`golint` doesn't version releases, so in order to get the latest version & force
a new docker image build, we have a `DATE_BUILT` file: to force an update of
golint, run `date > DATE_BUILT` locally & then rebuild the docker image.

### Need help?

For help with Golint, [check out their documentation](https://github.com/golang/lint).

If you're running into a Code Climate issue, first look over this project's [GitHub Issues](https://github.com/codeclimate/codeclimate-golint/issues), as your question may have already been covered. If not, [go ahead and open a support ticket with us](https://codeclimate.com/help).
