package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vektah/gqlparser/ast"
	"regexp"
)

func ReWriteQl(gql *ast.QueryDocument) error {
	if len(gql.Operations) == 0 || len(gql.Operations[0].SelectionSet) == 0 {
		logrus.Warnf("[plugin][regexp_query_selection]Operations or selections is empty")
		return nil
	}
	field, ok := gql.Operations[0].SelectionSet[0].(*ast.Field)
	if !ok {
		logrus.Warnf("[plugin][regexp_query_selection]section convert to ast.field error")
		return nil
	}
	logrus.Debugf("[plugin][regexp_query_selection]field name:%s ,alias: %s", field.Name, field.Alias)
	regexps := viper.GetStringMapString(`rewriteql_regexps`)
	for k, v := range regexps {
		logrus.Debugf(`[rewrite_query_selection]regexp:%s ,replace_regexp:%s`, k, v)
		compile := regexp.MustCompile(k)
		newName := compile.ReplaceAllString(field.Name, v)
		logrus.Debugf(`[rewrite_query_selection]field name:%s,new name:%s`, field.Name, newName)
		field.Name = newName
	}

	return nil
}

func main() {
	println("print plugin main")
}
