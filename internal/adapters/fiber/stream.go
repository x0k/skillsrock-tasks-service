package fiber_adapter

import (
	"bufio"
	"encoding/json"
	"iter"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

const chunkSize = 1024 * 1024 * 10

func JSONSequence[T any](c *fiber.Ctx, items iter.Seq2[T, error]) error {
	size := reflect.TypeFor[T]().Size()
	flushRate := int(max(chunkSize/size, 1))
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		if err := w.WriteByte('['); err != nil {
			return
		}
		var i int
		encoder := json.NewEncoder(w)
		for item, err := range items {
			if err != nil {
				return
			}
			if i > 0 {
				if err := w.WriteByte(','); err != nil {
					return
				}
			}
			if err := encoder.Encode(item); err != nil {
				return
			}
			if i%flushRate == 0 {
				w.Flush()
			}
		}
		if err := w.WriteByte(']'); err != nil {
			return
		}
		w.Flush()
	})
	return nil
}
