import React, { useState, useEffect, useRef } from 'react';
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  FlatList,
  StyleSheet,
  ActivityIndicator,
  Alert,
  ListRenderItem
} from 'react-native';
import { SafeAreaProvider, SafeAreaView } from 'react-native-safe-area-context';

interface MessageStats {
  ack: number;
  publish: number;
}

interface Queue {
  name: string;
  messages: number;
  consumers: number;
  message_stats?: MessageStats;
}

export default function App() {
  const [ipAddress, setIpAddress] = useState<string>('192.168.1.35:8080'); 
  const [queues, setQueues] = useState<Queue[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  
  // To verify if connection is established to start polling
  const [isConnected, setIsConnected] = useState<boolean>(false);

  // FETCH FUNCTION: Supports 'silent' mode to avoid spinner flickering
  const fetchQueues = async (silent: boolean = false) => {
    if (!silent) setLoading(true);
    
    let formattedIp = ipAddress;
    if (!formattedIp.startsWith('http')) {
      formattedIp = `http://${formattedIp}`;
    }
    
    const apiUrl = `${formattedIp}/queues`;

    try {
      const response = await fetch(apiUrl);
      if (!response.ok) {
         throw new Error(`HTTP Error: ${response.status}`);
      }
      const data: Queue[] = await response.json();
      setQueues(data);
      setIsConnected(true); // Connection successful, keep polling
    } catch (error) {
      console.error(error);
      if (!silent) {
        Alert.alert("Connection Error", "Could not reach server.");
      }
    } finally {
      if (!silent) setLoading(false);
    }
  };

  // AUTO-REFRESH LOGIC (POLLING)
  useEffect(() => {
    let interval: any;

    if (isConnected) {
      // If connected, fetch every 2 seconds silently
      interval = setInterval(() => {
        fetchQueues(true); // true = silent mode
      }, 2000);
    }

    // Cleanup: Stop timer when screen closes or connection stops
    return () => clearInterval(interval);
  }, [isConnected, ipAddress]);


  const renderQueueItem: ListRenderItem<Queue> = ({ item }) => {
    const isDlq = item.name.includes("dlq");
    const isCritical = isDlq || item.messages > 50;
    
    const cardBg = isCritical ? "#ffebee" : "#e3f2fd";
    const cardBorder = isCritical ? "#ef5350" : "#42a5f5";
    const textColor = isCritical ? '#c62828' : '#1565c0';

    return (
      <View style={[styles.card, { backgroundColor: cardBg, borderColor: cardBorder }]}>
        <View style={styles.cardHeader}>
          <Text style={[styles.queueName, { color: textColor }]}>
            {item.name}
          </Text>
          {isCritical && <Text style={styles.badge}>⚠️ CRITICAL</Text>}
        </View>

        <View style={styles.statsRow}>
          <View style={styles.statBox}>
            <Text style={styles.statLabel}>PENDING</Text>
            <Text style={[styles.statValue, { color: item.messages > 0 ? '#d32f2f' : '#333' }]}>
              {item.messages}
            </Text>
          </View>
          
          <View style={styles.statBox}>
            <Text style={styles.statLabel}>PROCESSED</Text>
            <Text style={[styles.statValue, { color: '#2e7d32' }]}>
              {item.message_stats?.ack ?? 0}
            </Text>
          </View>

          <View style={styles.statBox}>
            <Text style={styles.statLabel}>CONSUMERS</Text>
            <Text style={styles.statValue}>{item.consumers}</Text>
          </View>
        </View>
      </View>
    );
  };

  return (
    <SafeAreaProvider>
      <SafeAreaView style={styles.container}>
        <Text style={styles.headerTitle}>🐰 RabbitMQ Monitor</Text>
        
        {/* Status Indicator */}
        {isConnected && (
            <View style={styles.liveIndicator}>
                <View style={styles.dot} />
                <Text style={styles.liveText}>LIVE UPDATING (2s)</Text>
            </View>
        )}

        <View style={styles.inputContainer}>
          <Text style={styles.label}>Server IP (e.g., 192.168.1.35:8080)</Text>
          <TextInput
            style={styles.input}
            value={ipAddress}
            onChangeText={(text) => {
                setIpAddress(text);
                setIsConnected(false); // Stop polling if IP changes
            }}
            placeholder="192.168.1.35:8080"
            keyboardType="url"
            autoCapitalize="none"
            autoCorrect={false}
          />
          <TouchableOpacity 
              style={[styles.button, isConnected ? styles.buttonActive : {}]}
              onPress={() => fetchQueues(false)}
              disabled={loading}
          >
            {loading ? (
              <ActivityIndicator color="white" />
            ) : (
              <Text style={styles.buttonText}>
                  {isConnected ? "CONNECTED (Tap to Restart)" : "CONNECT & START"}
              </Text>
            )}
          </TouchableOpacity>
        </View>

        <FlatList
          data={queues}
          keyExtractor={(item) => item.name}
          renderItem={renderQueueItem}
          contentContainerStyle={styles.listContent}
          ListEmptyComponent={
            !loading ? (
              <Text style={styles.emptyText}>
                No data. Enter IP address to connect.
              </Text>
            ) : null
          }
        />
      </SafeAreaView>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#f5f5f5' },
  headerTitle: { fontSize: 24, fontWeight: 'bold', textAlign: 'center', marginTop: 20, marginBottom: 5, color: '#333' },
  
  liveIndicator: { flexDirection: 'row', justifyContent:'center', alignItems:'center', marginBottom: 15 },
  dot: { width: 8, height: 8, borderRadius: 4, backgroundColor: '#4caf50', marginRight: 6 },
  liveText: { fontSize: 12, color: '#4caf50', fontWeight: 'bold' },

  inputContainer: { backgroundColor: 'white', padding: 15, marginHorizontal: 15, borderRadius: 10, elevation: 3, shadowColor: '#000', shadowOpacity: 0.1, shadowRadius: 5 },
  label: { fontSize: 12, color: '#666', marginBottom: 5, fontWeight:'600' },
  input: { borderWidth: 1, borderColor: '#ddd', borderRadius: 8, padding: 12, fontSize: 16, marginBottom: 12, backgroundColor: '#fafafa' },
  button: { backgroundColor: '#757575', padding: 14, borderRadius: 8, alignItems: 'center' },
  buttonActive: { backgroundColor: '#2196f3' }, // Blue when connected
  buttonText: { color: 'white', fontWeight: 'bold', fontSize: 16 },

  listContent: { padding: 15, paddingBottom: 50 },
  card: { borderWidth: 1, borderRadius: 12, padding: 16, marginBottom: 12, elevation: 2 },
  cardHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 15, borderBottomWidth: 1, borderBottomColor: 'rgba(0,0,0,0.05)', paddingBottom: 10 },
  queueName: { fontSize: 18, fontWeight: 'bold' },
  badge: { fontSize: 10, backgroundColor: '#ffebee', color: '#c62828', paddingHorizontal: 6, paddingVertical: 2, borderRadius: 4, fontWeight: 'bold', overflow: 'hidden' },
  
  statsRow: { flexDirection: 'row', justifyContent: 'space-between' },
  statBox: { alignItems: 'center', flex: 1 },
  statLabel: { fontSize: 10, color: '#666', fontWeight: 'bold', marginBottom: 2 },
  statValue: { fontSize: 18, fontWeight: '700' },
  
  emptyText: { textAlign: 'center', marginTop: 50, color: '#999', fontSize: 16 }
});