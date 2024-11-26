package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// defining a rule structure. The action field would be "allow" or "block"
type Rule struct {
	Protocol string `json:"protocol"`
	SourceIP string `json:"source_ip"`
	Port     string `json:"port"`
	Allow    bool   `json:"allow"`
}

var (
	sliceRules []Rule
	mutex      sync.Mutex
)

// add a rule
func AddRule(rule Rule) error {
	mutex.Lock()
	defer mutex.Unlock()

	if rule.SourceIP == "" || rule.Port == "" {
		return errors.New("invalid rule: missing source IP or port")
	}

	var rules []Rule
	rules_file, err := os.OpenFile("firewall_rules.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("could not open rules file: %v", err)
	}
	defer rules_file.Close()

	if stat, _ := rules_file.Stat(); stat.Size() > 0 {
		decoder := json.NewDecoder(rules_file)
		if err := decoder.Decode(&rules); err != nil {
			return fmt.Errorf("could not decode rules: %v", err)
		}
	}

	// check if rule already exists
	rules = append(rules, rule)

	if err := SaveRulesToFile(rules_file, rules); err != nil {
		return fmt.Errorf("could not save rules: %v", err)
	}

	sliceRules = append(sliceRules, rule)

	// // marshalling the struct to json
	// jsonBytes, err := json.Marshal(rule)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return err
	// }

	// // add the rule to file

	// // we need to put the pointer file at the end of the structures

	// _, err = rules_file.Write(jsonBytes)
	// if err != nil {
	// 	fmt.Println("Error adding to file: ", err)
	// 	return err
	// }

	return nil
}

func CheckRules(packet Packet) bool {
	mutex.Lock()
	defer mutex.Unlock()

	for _, rule := range sliceRules {
		if packet.SOURCE_IP == rule.SourceIP && packet.SOURCE_PORT == rule.Port {
			return rule.Allow
		}
	}

	return true
}

func SaveRulesToFile(file *os.File, rules []Rule) error {
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("could not truncate file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("could not seek file: %v", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(rules); err != nil {
		return fmt.Errorf("could not encode rules: %v", err)
	}

	return nil
}

func CheckFileRules(newRule Rule) (bool, error) {
	var rules []Rule
	rules_file, err := os.OpenFile("firewall_rules.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return true, fmt.Errorf("could not open rules file: %v", err)
	}
	defer rules_file.Close()

	if stat, _ := rules_file.Stat(); stat.Size() > 0 {
		decoder := json.NewDecoder(rules_file)
		if err := decoder.Decode(&rules); err != nil {
			return true, fmt.Errorf("could not decode rules: %v", err)
		}
	}

	// now the rules from the file are inside the rules slice.
	// we need to iterate
	for _, rule := range rules {
		if rule.Protocol == newRule.Protocol &&
			rule.SourceIP == newRule.SourceIP &&
			rule.Port == newRule.Port &&
			rule.Allow == newRule.Allow {
			return true, errors.New("rule already exists")
		}
	}

	return false, nil
}
