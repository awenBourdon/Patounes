import { View, Text, StyleSheet, ActivityIndicator } from 'react-native';
import { gql } from '@apollo/client';
import { useQuery } from '@apollo/client/react';

interface User {
  id: string;
  name: string;
  email: string;
}

interface UsersData {
  users: User[];
}

const USERS_QUERY = gql`
  query GetUsers {
    users {
      id
      name
      email
    }
  }
`;

export default function Index() {
  const { data, loading, error } = useQuery<UsersData>(USERS_QUERY);

  if (loading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#007AFF" />
        <Text style={styles.loading}>Chargement...</Text>
      </View>
    );
  }

  if (error) {
    return (
      <View style={styles.container}>
        <Text style={styles.error}>❌ Erreur:</Text>
        <Text style={styles.errorMessage}>{error.message}</Text>
      </View>
    );
  }

  const users = data?.users || [];

  return (
    <View style={styles.container}>
      <Text style={styles.title}>🐾 Patounes</Text>
      <Text style={styles.subtitle}>Utilisateurs GraphQL</Text>
      
      <View style={styles.userList}>
        {users.map((user) => ( // ← Plus besoin de any
          <View key={user.id} style={styles.userCard}>
            <Text style={styles.userName}>{user.name}</Text>
            <Text style={styles.userEmail}>{user.email}</Text>
          </View>
        ))}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
    padding: 20,
  },
  title: {
    fontSize: 32,
    fontWeight: 'bold',
    marginBottom: 10,
  },
  subtitle: {
    fontSize: 18,
    color: '#666',
    marginBottom: 30,
  },
  loading: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  error: {
    fontSize: 18,
    color: 'red',
    fontWeight: 'bold',
    marginBottom: 10,
  },
  errorMessage: {
    fontSize: 14,
    color: '#666',
    textAlign: 'center',
  },
  userList: {
    width: '100%',
  },
  userCard: {
    backgroundColor: '#fff',
    padding: 15,
    borderRadius: 8,
    marginBottom: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  userName: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  userEmail: {
    fontSize: 14,
    color: '#666',
  },
});