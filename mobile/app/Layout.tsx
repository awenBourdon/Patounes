
import { apolloClient } from '@/api/graphql';
import { ApolloProvider } from '@apollo/client/react';
import { Stack } from 'expo-router';

export default function RootLayout() {
  return (
    <ApolloProvider client={apolloClient}>
      <Stack />
    </ApolloProvider>
  );
}