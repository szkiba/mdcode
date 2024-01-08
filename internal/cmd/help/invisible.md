It is possible to use invisible code blocks. This is useful, for embedding test code or additional files in the markdown document. The invisible code block is also useful if you want to embed the entire file, but only want to display certain parts of it.

A markdown document can contain HTML elements. Unknown or unsupported HTML elements are usually not rendered by markdown renderers. Taking advantage of this, `mdcode` supports hiding code blocks using the standard `<script>` HTML element:

    <script type="text/markdown">
    ```js file=sample.js region=factorial

    ```
    </script>

Unfortunately, the GitHub markdown renderer renders the content of unsupported HTML elements as text. Therefore, `mdcode` also supports the use of a `<script>` element surrounded by an HTML comment to hide a code block.

    <!--<script type="text/markdown">
    ```js file=sample.js region=factorial

    ```
    </script>-->

*It is important to note that the opening character of the comment and the opening tag of the script element must be placed on the same line. Similarly, the closing tag of the script element and the closing tag of the comment must also be placed on the same line.*
