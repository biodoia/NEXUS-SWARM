package ternary

import (
  	"fmt"
  	"sync"
  	"time"

  	"github.com/google/uuid"
  )

// ════════════════════════════════════════════════════════════════
// TERNARY LOGIC ENGINE — Beyond Binary Thinking
// ════════════════════════════════════════════════════════════════
// Three-valued logic: TRUE | FALSE | UNKNOWN
// Inspired by Lukasiewicz logic, Kleene logic, and quantum superposition.
// Every decision in NEXUS passes through this engine.
// UNKNOWN is not absence of knowledge — it is a third state of being.

// Trit represents a ternary digit: the fundamental unit of ternary logic
type Trit int8

const (
  	FALSE   Trit = -1 // Definite negative
  	UNKNOWN Trit = 0  // Indeterminate / superposition
  	TRUE    Trit = 1  // Definite positive
  )

// String returns the CP437-styled representation
func (t Trit) String() string {
  	switch t {
      	case TRUE:
      		return "█ TRUE"
      	case FALSE:
      		return "░ FALSE"
      	case UNKNOWN:
      		return "▒ UNKNOWN"
      	default:
      		return "? INVALID"
      	}
  }

// Confidence returns a float64 confidence level
func (t Trit) Confidence() float64 {
  	switch t {
      	case TRUE:
      		return 1.0
      	case FALSE:
      		return 0.0
      	case UNKNOWN:
      		return 0.5
      	default:
      		return -1.0
      	}
  }

// TernaryResult holds a decision result with metadata
type TernaryResult struct {
  	ID         string    `json:"id"`
  	Value      Trit      `json:"value"`
  	Confidence float64   `json:"confidence"`
  	Reason     string    `json:"reason"`
  	Timestamp  time.Time `json:"timestamp"`
  	Depth      int       `json:"depth"` // recursive evaluation depth
  }

// Engine is the ternary logic evaluation engine
type Engine struct {
  	mu          sync.RWMutex
  	decisions   []TernaryResult
  	rules       map[string]TernaryRule
  	evalCount   uint64
  	truthTable  map[string]Trit
  }

// TernaryRule defines a named ternary evaluation rule
type TernaryRule struct {
  	Name     string
  	Evaluate func(inputs ...Trit) Trit
  	Weight   float64
  }

// NewEngine creates a new ternary logic engine
func NewEngine() *Engine {
  	e := &Engine{
      		decisions:  make([]TernaryResult, 0, 1024),
      		rules:      make(map[string]TernaryRule),
      		truthTable: make(map[string]Trit),
      	}
  	e.registerDefaultRules()
  	return e
  }

// registerDefaultRules sets up the fundamental ternary operations
func (e *Engine) registerDefaultRules() {
  	// Ternary AND (Kleene strong)
  	e.rules["AND"] = TernaryRule{
      		Name: "AND",
      		Evaluate: func(inputs ...Trit) Trit {
            			result := TRUE
            			for _, inp := range inputs {
                    				result = tritMin(result, inp)
                    			}
            			return result
            		},
      		Weight: 1.0,
      	}

  	// Ternary OR (Kleene strong)
  	e.rules["OR"] = TernaryRule{
      		Name: "OR",
      		Evaluate: func(inputs ...Trit) Trit {
            			result := FALSE
            			for _, inp := range inputs {
                    				result = tritMax(result, inp)
                    			}
            			return result
            		},
      		Weight: 1.0,
      	}

  	// Ternary NOT (Lukasiewicz)
  	e.rules["NOT"] = TernaryRule{
      		Name: "NOT",
      		Evaluate: func(inputs ...Trit) Trit {
            			if len(inputs) == 0 {
                    				return UNKNOWN
                    			}
            			return tritNeg(inputs[0])
            		},
      		Weight: 1.0,
      	}

  	// CONSENSUS — requires majority agreement
  	e.rules["CONSENSUS"] = TernaryRule{
      		Name: "CONSENSUS",
      		Evaluate: func(inputs ...Trit) Trit {
            			if len(inputs) == 0 {
                    				return UNKNOWN
                    			}
            			trueCount, falseCount, unknownCount := 0, 0, 0
            			for _, inp := range inputs {
                    				switch inp {
                              				case TRUE:
                              					trueCount++
                              				case FALSE:
                              					falseCount++
                              				case UNKNOWN:
                              					unknownCount++
                              				}
                    			}
            			total := len(inputs)
            			if trueCount > total/2 {
                    				return TRUE
                    			}
            			if falseCount > total/2 {
                    				return FALSE
                    			}
            			return UNKNOWN
            		},
      		Weight: 1.5,
      	}

  	// EVOLVE — biased toward action when uncertain
  	e.rules["EVOLVE"] = TernaryRule{
      		Name: "EVOLVE",
      		Evaluate: func(inputs ...Trit) Trit {
            			unknowns := 0
            			for _, inp := range inputs {
                    				if inp == UNKNOWN {
                              					unknowns++
                              				}
                    			}
            			// If more than 30% unknown, lean toward TRUE (action bias)
            			if float64(unknowns)/float64(len(inputs)) > 0.3 {
                    				return TRUE
                    			}
            			return e.rules["CONSENSUS"].Evaluate(inputs...)
            		},
      		Weight: 2.0,
      	}
  }

// Evaluate processes a decision through the ternary engine
func (e *Engine) Evaluate(ruleName string, inputs ...Trit) TernaryResult {
  	e.mu.Lock()
  	defer e.mu.Unlock()

  	e.evalCount++

  	rule, exists := e.rules[ruleName]
  	if !exists {
      		return TernaryResult{
            			ID:         uuid.New().String(),
            			Value:      UNKNOWN,
            			Confidence: 0.0,
            			Reason:     fmt.Sprintf("Rule '%s' not found", ruleName),
            			Timestamp:  time.Now(),
            		}
      	}

  	value := rule.Evaluate(inputs...)
  	confidence := value.Confidence() * rule.Weight
  	if confidence > 1.0 {
      		confidence = 1.0
      	}

  	result := TernaryResult{
      		ID:         uuid.New().String(),
      		Value:      value,
      		Confidence: confidence,
      		Reason:     fmt.Sprintf("Rule[%s] evaluated %d inputs", ruleName, len(inputs)),
      		Timestamp:  time.Now(),
      	}

  	e.decisions = append(e.decisions, result)
  	return result
  }

// AddRule registers a custom ternary rule
func (e *Engine) AddRule(name string, rule TernaryRule) {
  	e.mu.Lock()
  	defer e.mu.Unlock()
  	e.rules[name] = rule
  }

// Stats returns engine statistics
func (e *Engine) Stats() map[string]interface{} {
  	e.mu.RLock()
  	defer e.mu.RUnlock()
  	return map[string]interface{}{
      		"total_evaluations": e.evalCount,
      		"total_decisions":   len(e.decisions),
      		"registered_rules":  len(e.rules),
      	}
  }

// Helper functions for ternary arithmetic
func tritMin(a, b Trit) Trit {
  	if a < b {
      		return a
      	}
  	return b
  }

func tritMax(a, b Trit) Trit {
  	if a > b {
      		return a
      	}
  	return b
  }

func tritNeg(a Trit) Trit {
  	return -a
  }
