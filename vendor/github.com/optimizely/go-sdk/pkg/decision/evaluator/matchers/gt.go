/****************************************************************************
 * Copyright 2019-2020, Optimizely, Inc. and contributors                   *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

// Package matchers //
package matchers

import (
	"fmt"

	"github.com/optimizely/go-sdk/pkg/decision/evaluator/matchers/utils"
	"github.com/optimizely/go-sdk/pkg/entities"
	"github.com/optimizely/go-sdk/pkg/logging"
)

// GtMatcher matches against the "gt" match type
func GtMatcher(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (bool, error) {
	res, err := compare(condition, user, logger)
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

func compare(condition entities.Condition, user entities.UserContext, logger logging.OptimizelyLogProducer) (int, error) {
	if !user.CheckAttributeExists(condition.Name) {
		logger.Debug(fmt.Sprintf(logging.NullUserAttribute.String(), condition.StringRepresentation, condition.Name))
		return 0, fmt.Errorf(`no attribute named "%s"`, condition.Name)
	}

	if floatValue, ok := utils.ToFloat(condition.Value); ok {
		attributeValue, err := user.GetFloatAttribute(condition.Name)
		if err != nil {
			val, _ := user.GetAttribute(condition.Name)
			logger.Warning(fmt.Sprintf(logging.InvalidAttributeValueType.String(), condition.StringRepresentation, val, condition.Name))
			return 0, err
		}
		if floatValue < attributeValue {
			return 1, nil
		} else if floatValue > attributeValue {
			return -1, nil
		}
		return 0, nil
	}

	logger.Warning(fmt.Sprintf(logging.UnsupportedConditionValue.String(), condition.StringRepresentation))
	return 0, fmt.Errorf("audience condition %s evaluated to NULL because the condition value type is not supported", condition.Name)
}