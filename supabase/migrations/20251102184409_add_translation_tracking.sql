-- Add columns to track AI translation usage and language codes
ALTER TABLE flashcards
ADD COLUMN IF NOT EXISTS ai_translated_question BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS ai_translated_answer BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS question_lang VARCHAR(10),
ADD COLUMN IF NOT EXISTS answer_lang VARCHAR(10);

-- Add index for querying by translation status
CREATE INDEX IF NOT EXISTS idx_flashcards_ai_translated
ON flashcards(ai_translated_question, ai_translated_answer);

-- Add comments for clarity
COMMENT ON COLUMN flashcards.ai_translated_question IS 'Indicates if the question was AI-generated';
COMMENT ON COLUMN flashcards.ai_translated_answer IS 'Indicates if the answer was AI-generated';
COMMENT ON COLUMN flashcards.question_lang IS 'Language code for question (ISO 639-1: en, el, etc.)';
COMMENT ON COLUMN flashcards.answer_lang IS 'Language code for answer (ISO 639-1: en, el, etc.)';

