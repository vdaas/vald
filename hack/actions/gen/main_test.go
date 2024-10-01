// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

// NOT IMPLEMENTED BELOW
//
// func TestData_initPullRequestPaths(t *testing.T) {
// 	type fields struct {
// 		AliasImage        bool
// 		ConfigExists      bool
// 		Year              int
// 		ContainerType     ContainerType
// 		AppName           string
// 		BinDir            string
// 		BuildUser         string
// 		BuilderImage      string
// 		BuilderTag        string
// 		BuildStageName    string
// 		Maintainer        string
// 		PackageDir        string
// 		RootDir           string
// 		RuntimeImage      string
// 		RuntimeTag        string
// 		RuntimeUser       string
// 		Name              string
// 		BuildPlatforms    string
// 		Arguments         map[string]string
// 		Environments      map[string]string
// 		Entrypoints       []string
// 		EnvironmentsSlice []string
// 		ExtraCopies       []string
// 		ExtraImages       []string
// 		ExtraPackages     []string
// 		Preprocess        []string
// 		RunCommands       []string
// 		RunMounts         []string
// 		StageFiles        []string
// 		PullRequestPaths  []string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           AliasImage:false,
// 		           ConfigExists:false,
// 		           Year:0,
// 		           ContainerType:nil,
// 		           AppName:"",
// 		           BinDir:"",
// 		           BuildUser:"",
// 		           BuilderImage:"",
// 		           BuilderTag:"",
// 		           BuildStageName:"",
// 		           Maintainer:"",
// 		           PackageDir:"",
// 		           RootDir:"",
// 		           RuntimeImage:"",
// 		           RuntimeTag:"",
// 		           RuntimeUser:"",
// 		           Name:"",
// 		           BuildPlatforms:"",
// 		           Arguments:nil,
// 		           Environments:nil,
// 		           Entrypoints:nil,
// 		           EnvironmentsSlice:nil,
// 		           ExtraCopies:nil,
// 		           ExtraImages:nil,
// 		           ExtraPackages:nil,
// 		           Preprocess:nil,
// 		           RunCommands:nil,
// 		           RunMounts:nil,
// 		           StageFiles:nil,
// 		           PullRequestPaths:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           AliasImage:false,
// 		           ConfigExists:false,
// 		           Year:0,
// 		           ContainerType:nil,
// 		           AppName:"",
// 		           BinDir:"",
// 		           BuildUser:"",
// 		           BuilderImage:"",
// 		           BuilderTag:"",
// 		           BuildStageName:"",
// 		           Maintainer:"",
// 		           PackageDir:"",
// 		           RootDir:"",
// 		           RuntimeImage:"",
// 		           RuntimeTag:"",
// 		           RuntimeUser:"",
// 		           Name:"",
// 		           BuildPlatforms:"",
// 		           Arguments:nil,
// 		           Environments:nil,
// 		           Entrypoints:nil,
// 		           EnvironmentsSlice:nil,
// 		           ExtraCopies:nil,
// 		           ExtraImages:nil,
// 		           ExtraPackages:nil,
// 		           Preprocess:nil,
// 		           RunCommands:nil,
// 		           RunMounts:nil,
// 		           StageFiles:nil,
// 		           PullRequestPaths:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			data := &Data{
// 				AliasImage:        test.fields.AliasImage,
// 				ConfigExists:      test.fields.ConfigExists,
// 				Year:              test.fields.Year,
// 				ContainerType:     test.fields.ContainerType,
// 				AppName:           test.fields.AppName,
// 				BinDir:            test.fields.BinDir,
// 				BuildUser:         test.fields.BuildUser,
// 				BuilderImage:      test.fields.BuilderImage,
// 				BuilderTag:        test.fields.BuilderTag,
// 				BuildStageName:    test.fields.BuildStageName,
// 				Maintainer:        test.fields.Maintainer,
// 				PackageDir:        test.fields.PackageDir,
// 				RootDir:           test.fields.RootDir,
// 				RuntimeImage:      test.fields.RuntimeImage,
// 				RuntimeTag:        test.fields.RuntimeTag,
// 				RuntimeUser:       test.fields.RuntimeUser,
// 				Name:              test.fields.Name,
// 				BuildPlatforms:    test.fields.BuildPlatforms,
// 				Arguments:         test.fields.Arguments,
// 				Environments:      test.fields.Environments,
// 				Entrypoints:       test.fields.Entrypoints,
// 				EnvironmentsSlice: test.fields.EnvironmentsSlice,
// 				ExtraCopies:       test.fields.ExtraCopies,
// 				ExtraImages:       test.fields.ExtraImages,
// 				ExtraPackages:     test.fields.ExtraPackages,
// 				Preprocess:        test.fields.Preprocess,
// 				RunCommands:       test.fields.RunCommands,
// 				RunMounts:         test.fields.RunMounts,
// 				StageFiles:        test.fields.StageFiles,
// 				PullRequestPaths:  test.fields.PullRequestPaths,
// 			}
//
// 			data.initPullRequestPaths()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestData_initData(t *testing.T) {
// 	type fields struct {
// 		AliasImage        bool
// 		ConfigExists      bool
// 		Year              int
// 		ContainerType     ContainerType
// 		AppName           string
// 		BinDir            string
// 		BuildUser         string
// 		BuilderImage      string
// 		BuilderTag        string
// 		BuildStageName    string
// 		Maintainer        string
// 		PackageDir        string
// 		RootDir           string
// 		RuntimeImage      string
// 		RuntimeTag        string
// 		RuntimeUser       string
// 		Name              string
// 		BuildPlatforms    string
// 		Arguments         map[string]string
// 		Environments      map[string]string
// 		Entrypoints       []string
// 		EnvironmentsSlice []string
// 		ExtraCopies       []string
// 		ExtraImages       []string
// 		ExtraPackages     []string
// 		Preprocess        []string
// 		RunCommands       []string
// 		RunMounts         []string
// 		StageFiles        []string
// 		PullRequestPaths  []string
// 	}
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           AliasImage:false,
// 		           ConfigExists:false,
// 		           Year:0,
// 		           ContainerType:nil,
// 		           AppName:"",
// 		           BinDir:"",
// 		           BuildUser:"",
// 		           BuilderImage:"",
// 		           BuilderTag:"",
// 		           BuildStageName:"",
// 		           Maintainer:"",
// 		           PackageDir:"",
// 		           RootDir:"",
// 		           RuntimeImage:"",
// 		           RuntimeTag:"",
// 		           RuntimeUser:"",
// 		           Name:"",
// 		           BuildPlatforms:"",
// 		           Arguments:nil,
// 		           Environments:nil,
// 		           Entrypoints:nil,
// 		           EnvironmentsSlice:nil,
// 		           ExtraCopies:nil,
// 		           ExtraImages:nil,
// 		           ExtraPackages:nil,
// 		           Preprocess:nil,
// 		           RunCommands:nil,
// 		           RunMounts:nil,
// 		           StageFiles:nil,
// 		           PullRequestPaths:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           AliasImage:false,
// 		           ConfigExists:false,
// 		           Year:0,
// 		           ContainerType:nil,
// 		           AppName:"",
// 		           BinDir:"",
// 		           BuildUser:"",
// 		           BuilderImage:"",
// 		           BuilderTag:"",
// 		           BuildStageName:"",
// 		           Maintainer:"",
// 		           PackageDir:"",
// 		           RootDir:"",
// 		           RuntimeImage:"",
// 		           RuntimeTag:"",
// 		           RuntimeUser:"",
// 		           Name:"",
// 		           BuildPlatforms:"",
// 		           Arguments:nil,
// 		           Environments:nil,
// 		           Entrypoints:nil,
// 		           EnvironmentsSlice:nil,
// 		           ExtraCopies:nil,
// 		           ExtraImages:nil,
// 		           ExtraPackages:nil,
// 		           Preprocess:nil,
// 		           RunCommands:nil,
// 		           RunMounts:nil,
// 		           StageFiles:nil,
// 		           PullRequestPaths:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			data := &Data{
// 				AliasImage:        test.fields.AliasImage,
// 				ConfigExists:      test.fields.ConfigExists,
// 				Year:              test.fields.Year,
// 				ContainerType:     test.fields.ContainerType,
// 				AppName:           test.fields.AppName,
// 				BinDir:            test.fields.BinDir,
// 				BuildUser:         test.fields.BuildUser,
// 				BuilderImage:      test.fields.BuilderImage,
// 				BuilderTag:        test.fields.BuilderTag,
// 				BuildStageName:    test.fields.BuildStageName,
// 				Maintainer:        test.fields.Maintainer,
// 				PackageDir:        test.fields.PackageDir,
// 				RootDir:           test.fields.RootDir,
// 				RuntimeImage:      test.fields.RuntimeImage,
// 				RuntimeTag:        test.fields.RuntimeTag,
// 				RuntimeUser:       test.fields.RuntimeUser,
// 				Name:              test.fields.Name,
// 				BuildPlatforms:    test.fields.BuildPlatforms,
// 				Arguments:         test.fields.Arguments,
// 				Environments:      test.fields.Environments,
// 				Entrypoints:       test.fields.Entrypoints,
// 				EnvironmentsSlice: test.fields.EnvironmentsSlice,
// 				ExtraCopies:       test.fields.ExtraCopies,
// 				ExtraImages:       test.fields.ExtraImages,
// 				ExtraPackages:     test.fields.ExtraPackages,
// 				Preprocess:        test.fields.Preprocess,
// 				RunCommands:       test.fields.RunCommands,
// 				RunMounts:         test.fields.RunMounts,
// 				StageFiles:        test.fields.StageFiles,
// 				PullRequestPaths:  test.fields.PullRequestPaths,
// 			}
//
// 			data.initData()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestData_generateWorkflowStruct(t *testing.T) {
// 	type fields struct {
// 		AliasImage        bool
// 		ConfigExists      bool
// 		Year              int
// 		ContainerType     ContainerType
// 		AppName           string
// 		BinDir            string
// 		BuildUser         string
// 		BuilderImage      string
// 		BuilderTag        string
// 		BuildStageName    string
// 		Maintainer        string
// 		PackageDir        string
// 		RootDir           string
// 		RuntimeImage      string
// 		RuntimeTag        string
// 		RuntimeUser       string
// 		Name              string
// 		BuildPlatforms    string
// 		Arguments         map[string]string
// 		Environments      map[string]string
// 		Entrypoints       []string
// 		EnvironmentsSlice []string
// 		ExtraCopies       []string
// 		ExtraImages       []string
// 		ExtraPackages     []string
// 		Preprocess        []string
// 		RunCommands       []string
// 		RunMounts         []string
// 		StageFiles        []string
// 		PullRequestPaths  []string
// 	}
// 	type want struct {
// 		want *Workflow
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *Workflow, error) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want, got *Workflow, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       fields: fields {
// 		           AliasImage:false,
// 		           ConfigExists:false,
// 		           Year:0,
// 		           ContainerType:nil,
// 		           AppName:"",
// 		           BinDir:"",
// 		           BuildUser:"",
// 		           BuilderImage:"",
// 		           BuilderTag:"",
// 		           BuildStageName:"",
// 		           Maintainer:"",
// 		           PackageDir:"",
// 		           RootDir:"",
// 		           RuntimeImage:"",
// 		           RuntimeTag:"",
// 		           RuntimeUser:"",
// 		           Name:"",
// 		           BuildPlatforms:"",
// 		           Arguments:nil,
// 		           Environments:nil,
// 		           Entrypoints:nil,
// 		           EnvironmentsSlice:nil,
// 		           ExtraCopies:nil,
// 		           ExtraImages:nil,
// 		           ExtraPackages:nil,
// 		           Preprocess:nil,
// 		           RunCommands:nil,
// 		           RunMounts:nil,
// 		           StageFiles:nil,
// 		           PullRequestPaths:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           fields: fields {
// 		           AliasImage:false,
// 		           ConfigExists:false,
// 		           Year:0,
// 		           ContainerType:nil,
// 		           AppName:"",
// 		           BinDir:"",
// 		           BuildUser:"",
// 		           BuilderImage:"",
// 		           BuilderTag:"",
// 		           BuildStageName:"",
// 		           Maintainer:"",
// 		           PackageDir:"",
// 		           RootDir:"",
// 		           RuntimeImage:"",
// 		           RuntimeTag:"",
// 		           RuntimeUser:"",
// 		           Name:"",
// 		           BuildPlatforms:"",
// 		           Arguments:nil,
// 		           Environments:nil,
// 		           Entrypoints:nil,
// 		           EnvironmentsSlice:nil,
// 		           ExtraCopies:nil,
// 		           ExtraImages:nil,
// 		           ExtraPackages:nil,
// 		           Preprocess:nil,
// 		           RunCommands:nil,
// 		           RunMounts:nil,
// 		           StageFiles:nil,
// 		           PullRequestPaths:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			data := &Data{
// 				AliasImage:        test.fields.AliasImage,
// 				ConfigExists:      test.fields.ConfigExists,
// 				Year:              test.fields.Year,
// 				ContainerType:     test.fields.ContainerType,
// 				AppName:           test.fields.AppName,
// 				BinDir:            test.fields.BinDir,
// 				BuildUser:         test.fields.BuildUser,
// 				BuilderImage:      test.fields.BuilderImage,
// 				BuilderTag:        test.fields.BuilderTag,
// 				BuildStageName:    test.fields.BuildStageName,
// 				Maintainer:        test.fields.Maintainer,
// 				PackageDir:        test.fields.PackageDir,
// 				RootDir:           test.fields.RootDir,
// 				RuntimeImage:      test.fields.RuntimeImage,
// 				RuntimeTag:        test.fields.RuntimeTag,
// 				RuntimeUser:       test.fields.RuntimeUser,
// 				Name:              test.fields.Name,
// 				BuildPlatforms:    test.fields.BuildPlatforms,
// 				Arguments:         test.fields.Arguments,
// 				Environments:      test.fields.Environments,
// 				Entrypoints:       test.fields.Entrypoints,
// 				EnvironmentsSlice: test.fields.EnvironmentsSlice,
// 				ExtraCopies:       test.fields.ExtraCopies,
// 				ExtraImages:       test.fields.ExtraImages,
// 				ExtraPackages:     test.fields.ExtraPackages,
// 				Preprocess:        test.fields.Preprocess,
// 				RunCommands:       test.fields.RunCommands,
// 				RunMounts:         test.fields.RunMounts,
// 				StageFiles:        test.fields.StageFiles,
// 				PullRequestPaths:  test.fields.PullRequestPaths,
// 			}
//
// 			got, err := data.generateWorkflowStruct()
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_main(t *testing.T) {
// 	type want struct{}
// 	type test struct {
// 		name       string
// 		want       want
// 		checkFunc  func(want) error
// 		beforeFunc func(*testing.T)
// 		afterFunc  func(*testing.T)
// 	}
// 	defaultCheckFunc := func(w want) error {
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T,) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T,) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			main()
// 			if err := checkFunc(test.want); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
