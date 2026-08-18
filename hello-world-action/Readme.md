# Go Greeting Action

A simple GitHub Action written in **Go** that greets a person and returns the current time.

## Usage

```yaml
- name: Greet
  id: greet
  uses: Preetu-Sharma/go-greeting-action@v1
  with:
    who-to-greet: "Preetu"

- name: Show time
  run: echo "Greeting time: ${{ steps.greet.outputs.time }}"
```

## Input

### `who-to-greet`

The name of the person to greet.

Default:

```text
World
```

## Output

### `time`

Returns the time when the greeting was generated.

## Example

```text
Hello, Preetu!
Greeting time: 2026-08-18T10:30:00Z
```

## Built With

* Go
* Docker
* GitHub Actions
