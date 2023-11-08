/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package config

import (
	"errors"
	"reflect"
)

// DeepCopy make a compele copy for a struct value
func DeepCopy(s, d any) error {
	sType := reflect.TypeOf(s)
	sValue := reflect.ValueOf(s)
	if sType.Kind() != reflect.Pointer {
		return errors.New("source object must be an Pointer")
	}
	sType = sType.Elem()

	dType := reflect.TypeOf(d)
	dValue := reflect.ValueOf(d)
	if dType.Kind() != reflect.Pointer {
		return errors.New("dest object must be an Pointer")
	}
	dType = dType.Elem()

	if sType != dType {
		return errors.New("source and dest object type is not match")
	}

	if sType.Kind() != reflect.Struct {
		return errors.New("object type is not struct")
	}

	return deepCopy(sValue, dValue)
}

func deepCopy(s, d reflect.Value) error {
	sValue := reflect.Indirect(s)
	kind := sValue.Kind()
	var err error
	switch kind {
	case reflect.Struct:
		err = deepCopyStruct(s, d)
	case reflect.String:
		err = deepCopyString(s, d)
	}

	return err
}

func deepCopyStruct(s, d reflect.Value) error {
	// 	var structs = make([]any, 1, 5)
	// 	structs[0] = d
	// 	type field struct {
	// 		field reflect.StructField
	// 		val   reflect.Value
	// 	}
	//
	// 	fields := []field{}
	//
	// for len(structs) > 0 {
	// 	structData := structs[0]
	// 	structs = structs[1:]
	dValue := reflect.Indirect(d)
	sValue := reflect.Indirect(s)
	dType := dValue.Type()

	for i := 0; i < dType.NumField(); i++ {
		dfieldValue := dValue.Field(i)
		sfieldValue := sValue.Field(i)
		deepCopy(sfieldValue, dfieldValue)
	}

	return nil
}

func deepCopyString(s, d reflect.Value) error {
	d.SetString(s.String())
	return nil
}
