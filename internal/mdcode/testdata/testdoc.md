# Test document

## entire file

```js file=entire.js
function add(a, b) {
    return a + b
}
```

<script type="text/markdown">
```js file=entire-script.js
function add(a, b) {
    return a + b
}
```
</script>


<!--<script type="text/markdown">
```js file=entire-comment.js
function add(a, b) {
    return a + b
}
```
</script>-->

## file region

<!--<script type="text/markdown">
```go file=partial.go outline=true
package main

// #region function
// #endregion
```
</script>-->

```go file=partial.go region=function
func add(a, b uint64) uint64 {
	return a + b
}

```
