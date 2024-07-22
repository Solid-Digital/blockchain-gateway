// Code generated by go-swagger; DO NOT EDIT.

package hello

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// HelloReader is a Reader for the Hello structure.
type HelloReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *HelloReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewHelloOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewHelloOK creates a HelloOK with default headers values
func NewHelloOK() *HelloOK {
	return &HelloOK{
		ContentType: "text/plain; charset=utf-8",
	}
}

/*HelloOK handles this case with default header values.


 */
type HelloOK struct {
	ContentType string

	Payload string
}

func (o *HelloOK) Error() string {
	return fmt.Sprintf("[GET /][%d] helloOK  %+v", 200, o.Payload)
}

func (o *HelloOK) GetPayload() string {
	return o.Payload
}

func (o *HelloOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response header Content-Type
	o.ContentType = response.GetHeader("Content-Type")

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
