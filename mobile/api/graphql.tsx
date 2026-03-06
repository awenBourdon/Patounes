import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';

const API_URL = process.env.EXPO_PUBLIC_API_URL;

const httpLink = new HttpLink({
    uri: `${API_URL}/graphql`,
    credentials: 'include',
});

export const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache(),
});