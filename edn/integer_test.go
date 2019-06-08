// Copyright 2018-2019 Workiva Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package edn

import (
	"github.com/Workiva/eva-client-go/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integer in EDN", func() {
	Context("", func() {

		It("should initialize without issue", func() {
			lexer, err := newLexer()
			Ω(err).Should(BeNil())

			delete(typeFactories, IntegerType)
			err = initInteger(lexer)
			Ω(err).Should(BeNil())
			_, has := typeFactories[IntegerType]
			Ω(has).Should(BeTrue())

			err = initInteger(lexer)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidFactory))
		})

		It("should create elements from the factory", func() {
			v := int64(123)

			elem, err := typeFactories[IntegerType](v)
			Ω(err).Should(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(IntegerType))
			Ω(elem.Value()).Should(BeEquivalentTo(v))
		})

		It("should not create elements from the factory if the input is not a the right type", func() {
			v := "foo"

			elem, err := typeFactories[IntegerType](v)
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrInvalidInput))
			Ω(elem).Should(BeNil())
		})

		It("should panic if the base factory errors.", func() {
			origFac := baseFactory
			baseFactory = func() elementFactory { return &breakerFactory{} }

			wrapper := func() {
				NewIntegerElement(123)
			}

			Ω(wrapper).Should(Panic())
			baseFactory = origFac
		})
	})

	Context("with the default marshaller", func() {

		testValue := int64(12345)

		It("should create an integer value with no error", func() {
			elem := NewIntegerElement(testValue)
			Ω(elem).ShouldNot(BeNil())
			Ω(elem.ElementType()).Should(BeEquivalentTo(IntegerType))
			Ω(elem.Value()).Should(BeEquivalentTo(testValue))
		})

		It("should serialize the integer without an issue", func() {
			elem := NewIntegerElement(testValue)
			Ω(elem).ShouldNot(BeNil())

			edn, err := elem.Serialize(EvaEdnMimeType)
			Ω(err).Should(BeNil())
			Ω(edn).Should(BeEquivalentTo("12345"))
		})

		It("should serialize the integer without an issue", func() {
			elem := NewIntegerElement(testValue)
			Ω(elem).ShouldNot(BeNil())

			_, err := elem.Serialize(SerializerMimeType("InvalidType"))
			Ω(err).ShouldNot(BeNil())
			Ω(err).Should(test.HaveMessage(ErrUnknownMimeType))
		})
	})

	Context("Parsing", func() {
		runParserTests(IntegerType,
			&testDefinition{"0", 0},
			&testDefinition{"+0", 0},
			&testDefinition{"-0", 0},
			&testDefinition{"1", 1},
			&testDefinition{"-1", -1},
			&testDefinition{"1234", 1234},
			&testDefinition{"0N", 0},
			&testDefinition{"+0N", 0},
			&testDefinition{"-0N", 0},
			&testDefinition{"1N", 1},
			&testDefinition{"-1N", -1},
			&testDefinition{"1234N", 1234},
		)
	})
})
