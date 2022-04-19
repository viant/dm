package xml

import (
	"fmt"
)

type span struct {
	start int
	end   int
}

func extractAttributes(offset int, input []byte) ([][2]span, error) {
	result := make([][2]span, 0)

	var currentOffset int
	for currentOffset < len(input) && !isWhitespace(input[currentOffset]) {
		currentOffset++
	}

	if currentOffset == len(input) {
		return [][2]span{}, nil
	}

	for currentOffset < len(input) {
		for currentOffset < len(input) && isWhitespace(input[currentOffset]) {
			currentOffset += 1
		}

		if currentOffset >= len(input) || input[currentOffset] == '/' || input[currentOffset] == '>' {
			return result, nil
		}

		spans, newOffset, err := matchAttribute(input[currentOffset:], currentOffset+offset)
		if err != nil {
			return nil, err
		}

		currentOffset = newOffset - offset + 1
		result = append(result, spans)
	}

	return result, nil
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\t' || b == '\r' || b == '\v' || b == '\f'
}

func matchAttribute(input []byte, offset int) ([2]span, int, error) {
	result := [2]span{}

	var i int
	var b byte
outer:
	for i, b = range input {
		switch b {
		case '=':
			result[0] = span{
				start: offset,
				end:   offset + i,
			}

			i++
			break outer
		case ' ', '\n', '\t', '\r', '\v', '\f', '>':
			return [2]span{
				{
					start: offset,
					end:   offset + i - 1,
				},
				{
					end:   offset + i - 1,
					start: offset + i - 1,
				},
			}, offset + i, nil
		}
	}

	foundQuote := -1
	for ; i < len(input); i++ {
		b = input[i]
		switch b {
		case '"', '\'':
			if foundQuote == -1 {
				foundQuote = i
			} else {
				if input[foundQuote] == b {
					result[1] = span{
						start: offset + foundQuote + 1,
						end:   offset + i,
					}

					return result, offset + i, nil
				}
			}
		case ' ', '\n', '\t', '\r', '\v', '\f':
			if foundQuote == -1 {
				result[1] = span{
					start: offset + 1,
					end:   offset + i - 1,
				}

				return result, offset + i, nil
			}
		}
	}

	return [2]span{}, -1, fmt.Errorf("not found attribute value %v", input)
}
