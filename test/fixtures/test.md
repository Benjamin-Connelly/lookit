# Markdown Test Document

This is a comprehensive test of markdown rendering capabilities.

## Headings

### Level 3 Heading
#### Level 4 Heading
##### Level 5 Heading
###### Level 6 Heading

## Text Formatting

This paragraph demonstrates **bold text**, *italic text*, and ***bold italic text***.

You can also use __bold__ and _italic_ with underscores.

## Lists

### Unordered List

- First item
- Second item
  - Nested item 1
  - Nested item 2
    - Deeply nested item
- Third item

### Ordered List

1. First step
2. Second step
   1. Sub-step A
   2. Sub-step B
3. Third step

### Task List

- [x] Completed task
- [ ] Incomplete task
- [ ] Another incomplete task

## Links

Check out [GitHub](https://github.com) for more information.

Visit [lookit](https://github.com/yourusername/lookit) on GitHub.

## Code

### Inline Code

Use `npm install` to install packages. The `createMarkdownTemplate()` function handles rendering.

### Code Blocks

JavaScript example:

```javascript
function greet(name) {
  console.log(`Hello, ${name}!`);
  return `Welcome to lookit`;
}

const result = greet('World');
```

Python example:

```python
def calculate_sum(a, b):
    """Calculate the sum of two numbers."""
    return a + b

result = calculate_sum(10, 20)
print(f"The result is: {result}")
```

Shell script:

```bash
#!/bin/bash

echo "Starting lookit server..."
npm start

# Check if server is running
if [ $? -eq 0 ]; then
    echo "Server started successfully!"
fi
```

## Tables

| Feature | Status | Priority |
|---------|--------|----------|
| Markdown rendering | ✅ Done | High |
| Code highlighting | ✅ Done | High |
| Directory listing | 🚧 In Progress | Medium |
| Binary files | 📋 Planned | Low |

## Blockquotes

> This is a blockquote.
> It can span multiple lines.
>
> And contain multiple paragraphs.

> **Note:** Blockquotes can contain other markdown elements like **bold** and *italic* text.

## Horizontal Rules

Above the line

---

Below the line

## Images

![Placeholder Image](https://via.placeholder.com/600x300)

## Mixed Content

Here's a complex example combining multiple elements:

1. **First Point**: This demonstrates a list with various formatting

   ```json
   {
     "name": "lookit",
     "version": "1.0.0",
     "description": "Beautiful file viewer"
   }
   ```

2. **Second Point**: With a nested blockquote

   > Important: Always test your markdown rendering!

3. **Third Point**: And a table

   | Column 1 | Column 2 |
   |----------|----------|
   | Value A  | Value B  |

## Special Characters

Testing special characters: `<script>alert('xss')</script>` should be escaped.

HTML entities: &copy; 2024 &mdash; Testing &amp; validation.

## Conclusion

This document tests various markdown features including:

- Headings (all levels)
- Text formatting (bold, italic, combined)
- Lists (ordered, unordered, nested, tasks)
- Code blocks with syntax highlighting
- Tables with alignment
- Blockquotes
- Horizontal rules
- Links and images
- Mixed content

If all these elements render correctly, the markdown handler is working properly!
