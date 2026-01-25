# GraphQL Testing Guide

This document contains example GraphQL queries and mutations to test the API.

## Mutations

### Create Account
```graphql
mutation CreateAccount {
  createAccount(input: {
    name: "John Doe"
  }) {
    id
    name
    orders {
      id
    }
  }
}
```

### Create Product
```graphql
mutation CreateProduct {
  createProduct(input: {
    name: "MacBook Pro"
    description: "M3 Max, 16-inch"
    price: 3499.99
  }) {
    id
    name
    description
    price
  }
}
```

### Create Order
Replace `YOUR_ACCOUNT_ID` and `YOUR_PRODUCT_ID` with actual IDs returned from the previous mutations.

```graphql
mutation CreateOrder {
  createOrder(input: {
    accountId: "YOUR_ACCOUNT_ID"
    products: [
      {
        id: "YOUR_PRODUCT_ID"
        quantity: 1
      }
    ]
  }) {
    id
    createdAt
    totalPrice
    products {
      name
      quantity
      price
    }
  }
}
```

## Queries

### List Accounts
```graphql
query ListAccounts {
  accounts(pagination: { skip: 0, take: 10 }) {
    id
    name
    orders {
      id
      totalPrice
    }
  }
}
```

### Get Account by ID
```graphql
query GetAccount {
  accounts(id: "YOUR_ACCOUNT_ID") {
    id
    name
    orders {
      id
      products {
        name
      }
    }
  }
}
```

### List Products
```graphql
query ListProducts {
  products(pagination: { skip: 0, take: 10 }) {
    id
    name
    price
  }
}
```

### Search Products
```graphql
query SearchProducts {
  products(query: "MacBook") {
    id
    name
    description
    price
  }
}
```
