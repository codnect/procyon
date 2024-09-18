package container

/*

type MockDefinitionRegistry struct {
	mock.Mock
}

func (r *MockDefinitionRegistry) Register(definition *Definition) error {
	got := r.Called(definition)

	if len(got) == 0 || got[0] == nil {
		return nil
	}

	return got[0].(error)
}

func (r *MockDefinitionRegistry) Remove(name string) error {
	got := r.Called(name)

	if len(got) == 0 || got[0] == nil {
		return nil
	}

	return got[0].(error)
}

func (r *MockDefinitionRegistry) Find(filters ...filter.Filter) (*Definition, error) {
	got := r.Called(filters)

	if len(got) == 0 {
		return nil, nil
	}

	if got[1] == nil {
		return got[0].(*Definition), nil
	}

	return got[0].(*Definition), got[1].(error)
}

func (r *MockDefinitionRegistry) FindFirst(filters ...filter.Filter) (*Definition, bool) {
	got := r.Called(filters)

	if len(got) == 0 {
		return nil, false
	}

	return got[0].(*Definition), got[1].(bool)
}

func (r *MockDefinitionRegistry) List(filters ...filter.Filter) []*Definition {
	got := r.Called(filters)
	return got[0].([]*Definition)
}

func (r *MockDefinitionRegistry) Contains(name string) bool {
	got := r.Called(name)
	return got[0].(bool)
}

func (r *MockDefinitionRegistry) Names() []string {
	got := r.Called()
	return got[0].([]string)
}

func (r *MockDefinitionRegistry) Count() int {
	got := r.Called()
	return got[0].(int)
}

func TestObjectDefinitionRegistry_RegisterShouldRegisterDefinitionSuccessfully(t *testing.T) {
	registry := newDefinitionRegistry()
	definition, err := MakeDefinition(anyConstructorFunction, Named("anyObjectName"))
	assert.Nil(t, err)
	err = registry.Register(definition)
	assert.Nil(t, err)

	assert.Contains(t, registry.definitionMap, "anyObjectName")
}

func TestObjectDefinitionRegistry_RegisterShouldReturnErrorIfDefinitionWithSameNameIsAlreadyRegistered(t *testing.T) {
	registry := newDefinitionRegistry()
	anyDefinition, err := MakeDefinition(anyConstructorFunction, Named("anyObjectName"))
	assert.Nil(t, err)
	err = registry.Register(anyDefinition)
	assert.Nil(t, err)

	err = registry.Register(anyDefinition)
	assert.Equal(t, "definition with name 'anyObjectName' already exists", err.Error())
}

func TestObjectDefinitionRegistry_ContainsShouldReturnTrueIfDefinitionExists(t *testing.T) {
	registry := newDefinitionRegistry()
	anyDefinition, err := MakeDefinition(anyConstructorFunction, Named("anyObjectName"))
	assert.Nil(t, err)
	err = registry.Register(anyDefinition)
	assert.Nil(t, err)

	assert.True(t, registry.Contains("anyObjectName"))
}

func TestObjectDefinitionRegistry_ContainsShouldReturnFalseIfDefinitionDoesNotExist(t *testing.T) {
	registry := newSingletonObjectRegistry()
	assert.False(t, registry.Contains("anyObjectName"))
}

func TestObjectDefinitionRegistry_Find(t *testing.T) {

	type args struct {
		filter []filter.Filter
	}

	type fields struct {
		definitions map[string]*Definition
	}

	anyDefinition, _ := MakeDefinition(anyConstructorFunction, Named("anyObjectName"))
	anotherDefinition, _ := MakeDefinition(anotherConstructorFunction, Named("anotherObjectName"))

	anyObject := &AnyType{}
	anyObjectType := reflect.TypeOf(anyObject)

	testCases := []struct {
		name    string
		fields  fields
		args    args
		want    *Definition
		wantErr string
	}{
		{
			name: "ShouldReturnNoFilterErrorWithoutFiltersIfThereIsNoAnyDefinition",
			fields: fields{
				definitions: map[string]*Definition{},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:    nil,
			wantErr: "at least one filter must be used",
		},
		{
			name: "ShouldReturnNoFilterErrorWithoutFiltersIfThereIsOnlyOneDefinition",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName": anyDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:    nil,
			wantErr: "at least one filter must be used",
		},
		{
			name: "ShouldReturnNoFilterErrorWithoutFiltersIfThereAreManyDefinitions",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:    nil,
			wantErr: "at least one filter must be used",
		},
		{
			name: "ShouldReturnDefinitionWithByNameFilterIfDefinitionWithNameExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want: anyDefinition,
		},
		{
			name:   "ShouldReturnErrorWithByNameFilterIfDefinitionWithNameDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want:    nil,
			wantErr: "not found definition with name 'anyObjectName'",
		},
		{
			name: "ShouldReturnDefinitionWithByTypeOfFilterIfObjectWithTypeExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: anyDefinition,
		},
		{
			name:   "ShouldReturnErrorWithByTypeOfFilterIfDefinitionWithTypeDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:    nil,
			wantErr: "not found definition with type '*component.AnyType'",
		},
		{
			name: "ShouldReturnDefinitionWithByTypeOfFilterIfThereIsOnlyOneObjectImplementingInterface",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnotherInterface](),
				},
			},
			want: anotherDefinition,
		},
		{
			name: "ShouldReturnErrorWithByTypeOfFilterIfThereIsMoreThanOneDefinitionImplementingInterface",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnyType](),
				},
			},
			want:    anyDefinition,
			wantErr: "cannot distinguish objects because too many matching found",
		},
		{
			name: "ShouldReturnDefinitionWithByTypeFilterIfObjectWithTypeExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want: anyDefinition,
		},
		{
			name:   "ShouldReturnErrorWithByTypeFilterIfDefinitionWithTypeDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want:    nil,
			wantErr: "not found definition with type '*component.AnyType'",
		},
		{
			name: "ShouldReturnDefinitionWithAllFiltersIfDefinitionExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want: anyDefinition,
		},
		{
			name:   "ShouldReturnErrorWithAllFiltersIfDefinitionDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:    nil,
			wantErr: "not found definition with name 'anyObjectName' and type '*component.AnyType'",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registry := newDefinitionRegistry()
			registry.definitionMap = testCase.fields.definitions

			got, err := registry.Find(testCase.args.filter...)
			if err != nil || testCase.wantErr != "" {
				if err != nil {
					assert.Equal(t, testCase.wantErr, err.Error(), "failed for test case '%s'", testCase.name)
				} else {
					assert.Nil(t, err, "want err '%s' but got nil", testCase.wantErr)
				}
			}

			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestObjectDefinitionRegistry_FindFirst(t *testing.T) {

	type args struct {
		filter []filter.Filter
	}

	type fields struct {
		definitions map[string]*Definition
	}

	anyDefinition, _ := MakeDefinition(anyConstructorFunction, Named("anyObjectName"))
	anotherDefinition, _ := MakeDefinition(anotherConstructorFunction, Named("anotherObjectName"))

	anyObject := &AnyType{}
	anyObjectType := reflect.TypeOf(anyObject)

	testCases := []struct {
		name     string
		fields   fields
		args     args
		want     *Definition
		wantIn   []*Definition
		wantBool bool
	}{
		{
			name: "ShouldReturnDefinitionWithoutFiltersIfThereIsOnlyOneDefinition",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName": anyDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:     anyDefinition,
			wantBool: true,
		},
		{
			name: "ShouldReturnDefinitionsWithoutFiltersIfThereAreManyDefinitions",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{},
			},
			want:     nil,
			wantIn:   []*Definition{anyDefinition, anotherDefinition},
			wantBool: true,
		},
		{
			name: "ShouldReturnDefinitionWithByNameFilterIfDefinitionWithNameExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want:     anyDefinition,
			wantBool: true,
		},
		{
			name:   "ShouldReturnErrorWithByNameFilterIfDefinitionWithNameDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
				},
			},
			want:     nil,
			wantBool: false,
		},
		{
			name: "ShouldReturnDefinitionWithByTypeOfFilterIfObjectWithTypeExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     anyDefinition,
			wantBool: true,
		},
		{
			name:   "ShouldReturnErrorWithByTypeOfFilterIfDefinitionWithTypeDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     nil,
			wantBool: false,
		},
		{
			name: "ShouldReturnDefinitionWithByTypeOfFilterIfThereIsOnlyOneObjectImplementingInterface",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnotherInterface](),
				},
			},
			want:     anotherDefinition,
			wantBool: true,
		},
		{
			name: "ShouldReturnDefinitionsWithByTypeOfFilterIfThereIsMoreThanOneDefinitionImplementingInterface",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByTypeOf[AnyType](),
				},
			},
			want:     nil,
			wantIn:   []*Definition{anyDefinition, anotherDefinition},
			wantBool: true,
		},
		{
			name: "ShouldReturnDefinitionWithByTypeFilterIfObjectWithTypeExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want:     anyDefinition,
			wantBool: true,
		},
		{
			name:   "ShouldReturnErrorWithByTypeFilterIfDefinitionWithTypeDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByType(anyObjectType),
				},
			},
			want:     nil,
			wantBool: false,
		},
		{
			name: "ShouldReturnDefinitionWithAllFiltersIfDefinitionExists",
			fields: fields{
				definitions: map[string]*Definition{
					"anyObjectName":     anyDefinition,
					"anotherObjectName": anotherDefinition,
				},
			},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     anyDefinition,
			wantBool: true,
		},
		{
			name:   "ShouldReturnErrorWithAllFiltersIfDefinitionDoesNotExist",
			fields: fields{},
			args: args{
				filter: []filter.Filter{
					filter.ByName("anyObjectName"),
					filter.ByType(anyObjectType),
					filter.ByTypeOf[*AnyType](),
				},
			},
			want:     nil,
			wantBool: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			registry := newDefinitionRegistry()
			registry.definitionMap = testCase.fields.definitions

			got, ok := registry.FindFirst(testCase.args.filter...)
			assert.Equal(t, testCase.wantBool, ok)
			if len(testCase.wantIn) != 0 {
				found := false
				for _, want := range testCase.wantIn {
					if want == got {
						found = true
						break
					}
				}

				assert.True(t, found, "not found any matching definition in wantIn %v", testCase.want)
			} else {
				assert.Equal(t, testCase.want, got)
			}
		})
	}
}

func TestObjectDefinitionRegistry_RemoveShouldDeleteDefinitionFromRegistryIfDefinitionWithNameExists(t *testing.T) {
	registry := newDefinitionRegistry()
	anyDefinition, err := MakeDefinition(anyConstructorFunction, Named("anyObjectName"))
	assert.Nil(t, err)
	err = registry.Register(anyDefinition)
	assert.Nil(t, err)

	err = registry.Remove("anyObjectName")
	assert.Nil(t, err)

	assert.Zero(t, len(registry.definitionMap))
	assert.Zero(t, registry.Count())
}

func TestObjectDefinitionRegistry_RemoveShouldReturnErrorIfDefinitionWithNameDoesNotExist(t *testing.T) {
	registry := newDefinitionRegistry()

	err := registry.Remove("anyObjectName")
	assert.Equal(t, "no found definition with name 'anyObjectName'", err.Error())
}

func TestObjectDefinitionRegistry_CountShouldReturnCountOfDefinitions(t *testing.T) {
	registry := newDefinitionRegistry()
	anyDefinition, err := MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = registry.Register(anyDefinition)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(registry.definitionMap))
	assert.Equal(t, 1, registry.Count())
}

func TestObjectDefinitionRegistry_NamesShouldReturnListOfObjectDefinitionNames(t *testing.T) {
	registry := newDefinitionRegistry()
	anyDefinition, err := MakeDefinition(anyConstructorFunction)
	assert.Nil(t, err)
	err = registry.Register(anyDefinition)
	assert.Nil(t, err)

	names := registry.Names()
	assert.Len(t, names, 1)
	assert.Contains(t, names, "anyType")
}
*/
