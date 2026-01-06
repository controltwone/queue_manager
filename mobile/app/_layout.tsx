import { Stack } from 'expo-router';

export default function Layout() {
  return (
    <Stack screenOptions={{
      headerStyle: { backgroundColor: '#f5f5f5' },
      headerShadowVisible: false,
      headerTitle: "", 
    }}>
      
      <Stack.Screen 
        name="index" 
        options={{ 
          headerShown: false 
        }} 
      />
    </Stack>
  );
}