# Supported Languages for Translation

## Overview
The flashcard application uses ISO 639-1 language codes for translation. OpenAI's translation API supports a wide range of languages.

## ISO 639-1 Language Codes
The backend accepts 2-letter lowercase language codes following the ISO 639-1 standard.

## Backend Validation
- **Format:** Exactly 2 lowercase letters (e.g., "en", "el", "fr")
- **No hardcoded whitelist:** Any valid ISO 639-1 code is accepted
- **OpenAI handles the rest:** The AI service will return an error if a language is not supported

## Common Languages for Frontend Dropdown

Here's a curated list of commonly used languages that should be included in the frontend dropdown:

### European Languages
- **en** - English
- **el** - Greek (Ελληνικά)
- **fr** - French (Français)
- **de** - German (Deutsch)
- **es** - Spanish (Español)
- **it** - Italian (Italiano)
- **pt** - Portuguese (Português)
- **nl** - Dutch (Nederlands)
- **pl** - Polish (Polski)
- **ru** - Russian (Русский)
- **uk** - Ukrainian (Українська)
- **sv** - Swedish (Svenska)
- **da** - Danish (Dansk)
- **no** - Norwegian (Norsk)
- **fi** - Finnish (Suomi)
- **cs** - Czech (Čeština)
- **ro** - Romanian (Română)
- **bg** - Bulgarian (Български)
- **hr** - Croatian (Hrvatski)
- **sr** - Serbian (Српски)
- **sk** - Slovak (Slovenčina)
- **sl** - Slovenian (Slovenščina)
- **et** - Estonian (Eesti)
- **lv** - Latvian (Latviešu)
- **lt** - Lithuanian (Lietuvių)

### Asian Languages
- **zh** - Chinese (中文)
- **ja** - Japanese (日本語)
- **ko** - Korean (한국어)
- **hi** - Hindi (हिन्दी)
- **th** - Thai (ไทย)
- **vi** - Vietnamese (Tiếng Việt)
- **id** - Indonesian (Bahasa Indonesia)
- **ms** - Malay (Bahasa Melayu)
- **tl** - Tagalog (Tagalog)

### Middle Eastern Languages
- **ar** - Arabic (العربية)
- **he** - Hebrew (עברית)
- **fa** - Persian (فارسی)
- **tr** - Turkish (Türkçe)

### Other Languages
- **sw** - Swahili (Kiswahili)
- **af** - Afrikaans
- **is** - Icelandic (Íslenska)
- **ga** - Irish (Gaeilge)
- **cy** - Welsh (Cymraeg)
- **eu** - Basque (Euskara)
- **ca** - Catalan (Català)

## Frontend Implementation Example

```json
// Example dropdown data structure
[
  { "code": "en", "name": "English", "native": "English" },
  { "code": "el", "name": "Greek", "native": "Ελληνικά" },
  { "code": "fr", "name": "French", "native": "Français" },
  { "code": "de", "name": "German", "native": "Deutsch" },
  { "code": "es", "name": "Spanish", "native": "Español" }
  // ... add more as needed
]
```

## API Request Format

When creating a translation request, use the language codes:

```json
POST /flashcards/translate
{
  "content": "hello",
  "from_lang": "en",
  "to_lang": "el"
}
```

## Notes for Frontend Team

1. **Start Simple:** Begin with 5-10 most common languages
2. **Show Native Names:** Display language names in their native scripts for better UX
3. **Default Languages:** Consider making "en" and "el" the default source/target
4. **Add More Later:** Easy to add more languages to dropdown without backend changes
5. **Error Handling:** If OpenAI doesn't support a rare language, backend will return an error

## Adding New Languages

To add a new language to the frontend dropdown:
1. Find the ISO 639-1 code: https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
2. Add to your frontend language list
3. No backend changes needed!

## Limitations

- Only 2-letter ISO 639-1 codes are supported (not 3-letter ISO 639-3)
- Language variants (e.g., en-US vs en-GB) are not differentiated
- OpenAI determines the language variant automatically based on context

