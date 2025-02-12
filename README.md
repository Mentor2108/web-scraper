# Scraper API Documentation

## Overview

This document provides details on how to send a request to the scraper API, including the expected request format, parameter descriptions, and example usage.

Currently this library only supports <u>chromedp in scrape-phase</u> and <u>goquery in process-phase</u>.

## Request Format

The scraper accepts JSON-based HTTP requests with the following structure:

```json
{
    "url": "<target-website-url>",
    "config": {
        "scrape_phase": {
            "library": "<scraping-library>",
            "wait_for": { //NOT YET IMPLEMENTED
                "duration": "<wait-time-in-seconds>",
                "selector": "<selector-for-to-wait>"
            }
        },
        "process_phase": {
            "library": "<processing-library>"
        },
        "root": "<CSS-selector-for-root-element>",
        "continueonerror": "<boolean>", // NOT YET IMPLEMENTED
        "depth": "<integer-depth>",
        "maxlimit": "<integer-max-elements>",
        "scrape_images": "<boolean>",
        "content": [
            {
                "name": "<identifier>",
                "type": "<text|section|table>",
                "selector": "<CSS-selector>",
                "text": {
                    "prefix": "<optional-prefix>",
                    "suffix": "<optional-suffix>"
                },
                "section": {
                    "prefix": "<optional-prefix>",
                    "suffix": "<optional-suffix>",
                    "start": "<start-selector>",
                    "end": "<end-selector>",
                    "title": ["<CSS-selectors-for-title>"] ,
                    "data": ["<CSS-selectors-for-data>"]
                },
                "table": {
                    "column_map": {
                        "key": "<CSS-selector-for-column-header>",
                        "value": "<CSS-selector-for-column-data>"
                    }
                }
            }
        ]
    }
}
```

## API Endpoints

- `POST /scraper/scrape/url/start` - Initiates scraping for a URL.

- `GET /scraper/list/scrapeddata/url?id=<JOB_ID>` - Retrieves scraped data for a job.

- `GET /content/file/id/:id` - Retrieves a file by its ID.

## Parameter Descriptions

### Top-Level Parameters

- `url` *(string, required)*: The target website URL to scrape.
- `config` *(object, required)*: Configuration for scraping behavior.

### Config Parameters

- `scrape_phase` *(object, required)*: Defines the scraping library and waiting mechanism.
  - `library` *(string, required)*: The library used for scraping (`chromedp`).
- `process_phase` *(object, optional)*: Defines the processing library.
  - `library` *(string, required)*: The library used for processing (`goquery`).
- `root` *(string, required, default="main")*: CSS selector defining the root of the scraping area.
- `depth` *(integer, optional, default=0)*: How deep the scraper should traverse the child links.
- `maxlimit` *(integer, optional, default=0)*: Maximum number of links to scrape when using depth `(Negative means all)`.
- `scrape_images` *(boolean, optional, default=false)*: Whether to scrape images.
- `content` *(array, required)*: List of elements to scrape for raw_text, each with a defined type and selector.

### Content Object

Each object inside `content` represents a specific scraping rule:

#### Common Fields

- `name` *(string, required)*: Identifier for the scraped data.
- `type` *(string, required)*: Type of data to extract (`text`, `section`, `table`).
- `selector` *(string, required)*: CSS selector of the target element.

#### Text Extraction (`type: "text"`)

- `text` *(object, optional)*:
  - `prefix` *(string, optional)*: String to prepend to extracted text.
  - `suffix` *(string, optional)*: String to append to extracted text.

#### Section Extraction (`type: "section"`)

Section Selector scrapes through the next siblings of `start` (if present) till the `end` selector is found (if present) or scrapes through the children of common-selector is `start` is not present.

- `section` *(object, required)*:
  - `prefix` *(string, optional)*: Prefix for extracted text.
  - `suffix` *(string, optional)*: Suffix for extracted text.
  - `start` *(string, option)*: CSS selector defining where the section starts. If not present, will scrape the first level children of common-selector
  - `end` *(string, optional)*: CSS selector defining where the section ends.
  - `title` *(array of strings, required, default:['h3'])*: List of CSS selectors for section titles.
  - `data` *(array of strings, required, default:['p'])*: List of CSS selectors for section data.

**NOTE: The selectors provided in `title` and `data` should directly match the list of sibling being scraped else they will be skipped.**

Eg:

```html
<head></head>
<body>
    <div class="main">
        <h2>Heading 1</h2>
        <p>Paragraph 1</p>
        <p>Paragraph 2</p>
        <h2>Heading 2</h2>
        <div>
            <p>Paragraph 3</p>
        </div>
    </div>
</body>
```

For the above HTML, the following content section will only give the following response: 
```json
"content" [
    "name": "scrape-section",
    "type": "section",
    "selector": ".main",
    "section": {
        "title": [
            "h2"
        ],
        "data": [
            "p"
        ]
    }
]
```

```md
SCRAPING_RESPONSE

## Heading 1
Paragraph 1
Paragraph 2
## Heading 2
```

This is because `data` key only has "p" as a selector, but the base level element is a "div" which is not recognised and hence will not be scraped. 

#### Table Extraction (`type: "table"`)

This is primarily used to scrape tables.  
There's two ways you can achieve this:  
1. By providing a `column_map` along with key/value pair of what value should correspond to what key in `each <trow> element`.
2. By providing a `column_names` list of strings. It will allot first child of `<trow>` to first column, 2nd child to 2nd column and so on. 

Keys:
- `table` *(object, required)*:
  - `column_map` *(object, required)*:
    - `key` *(string, required)*: CSS selector for table headers.
    - `value` *(string, required)*: CSS selector for table data.

## Example Request

```json
{
    "url": "https://en.wikipedia.org/wiki/India",
    "config": {
        "scrape_phase": {
            "library": "chromedp",
            "wait_for": {
                "duration": 10
            }
        },
        "process_phase": {
            "library": "goquery"
        },
        "root": "body",
        "depth": 2,
        "maxlimit": 4,
        "content": [
            {
                "name": "scrape-section",
                "type": "section",
                "selector": ".mw-content-ltr",
                "section": {
                    "suffix": "random_suffix_section\n",
                    "title": [
                        "div.mw-heading"
                    ],
                    "data": [
                        "p",
                        "li"
                    ],
                    "start": "p",
                    "end": "div:nth-child(16)"
                }
            },
            {
                "name": "scrape-table",
                "type": "table",
                "selector": "table",
                "table": {
                    "column_map": {
                        "key": "th",
                        "value": "td"
                    }
                }
            }
        ]
    }
}
```

## Notes

- The scraper can be configured with different libraries (`chromedp`, `goquery`, etc.) for different phases.
- The `scrape_phase` returns the raw HTML scraped over the root element.
- The `process_phase` determines how the extracted data is processed.
- `content` defines the data to be scraped, supporting text, sections, and tables.

## Contact

For further inquiries, please contact me directly.