package dal

import (
	"context"
	"testing"
)

//func TestWithPassword(t *testing.T) {
//	const password = "test-pwd"
//	txOptions := NewTransactionOptions(WithPassword(password))
//	if txOptions.Password() != password {
//		t.Errorf("unexpected password")
//	}
//}

func TestWithReadonly(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		txOptions := NewTransactionOptions(TxWithReadonly())
		if !txOptions.IsReadonly() {
			t.Errorf("expected to be readonly")
		}
	})
	t.Run("false", func(t *testing.T) {
		txOptions := NewTransactionOptions()
		if txOptions.IsReadonly() {
			t.Errorf("expected to be readonly")
		}
	})
}

func TestWithCrossGroup(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		txOptions := NewTransactionOptions(TxWithCrossGroup())
		if !txOptions.IsCrossGroup() {
			t.Errorf("expected to be true")
		}
	})
	t.Run("false", func(t *testing.T) {
		txOptions := NewTransactionOptions()
		if txOptions.IsCrossGroup() {
			t.Errorf("expected to be false")
		}
	})
}

func TestNewTransactionOptions(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		options := NewTransactionOptions()
		if options.IsReadonly() {
			t.Errorf("expected to be not readonly")
		}
		//if password := options.Password(); password != "" {
		//	t.Error("expected not to have a password, got: "+password)
		//}
	})
	t.Run("readonly", func(t *testing.T) {
		options := NewTransactionOptions(TxWithReadonly())
		if !options.IsReadonly() {
			t.Errorf("expected to be readonly")
		}
		//if password := options.Password(); password != "" {
		//	t.Errorf("expected not to have a password, got: %v", password)
		//}
	})
	//t.Run("password", func(t *testing.T) {
	//	const expectedPassword = "test-pwd"
	//	options := NewTransactionOptions(WithPassword(expectedPassword))
	//	if options.IsReadonly() {
	//		t.Errorf("expected not to be readonly")
	//	}
	//	if password := options.Password(); password != expectedPassword {
	//		t.Errorf("expected not to have a password equal to %v, got: %v", expectedPassword, password)
	//	}
	//})
}

type mockTx struct {
	options TransactionOptions
}

func (t mockTx) Options() TransactionOptions {
	return t.options
}

func TestGetTransaction(t *testing.T) {
	expected := mockTx{options: NewTransactionOptions()}
	ctx := context.Background()
	txCtx := NewContextWithTransaction(ctx, expected)
	actual := GetTransaction(txCtx)
	if actual != expected {
		t.Errorf("transactional context does not provide transaction's value")
	}
}

func TestGetNonTransactionalContext(t *testing.T) {
	expected := mockTx{options: NewTransactionOptions()}
	ctx := context.Background()
	txCtx := NewContextWithTransaction(ctx, expected)
	actual := GetNonTransactionalContext(txCtx)
	if actual != ctx {
		t.Errorf("transactional context does not provide original context")
	}
}
