vcl 4.0;

# Backend tanımı (Golang API Sunucusu)
backend default {
    .host = "host.docker.internal"; # Docker konteyner adı veya IP adresi
    .port = "8080";       # Golang API'nin çalıştığı port
}

# Gelen istekleri yöneten kısım
sub vcl_recv {
    if (req.method != "GET") {
        return (pass);  # GET dışında cacheleme yapma
    }

    return (hash);  # Cache kontrolü
}

# Cache davranışlarını kontrol eden bölüm
sub vcl_backend_response {
    # Sadece HTTP GET isteklerini cachele
    if (bereq.method == "GET") {
        set beresp.ttl = 60s;  # Cache süresi: 60 saniye
        set beresp.grace = 30s; # Cache'den yanıt verilemediğinde 30 saniye daha bekle
    } else {
        # Diğer HTTP metodlarını cacheleme
        set beresp.uncacheable = true;
    }
} 

# Cache'de bulunan bir yanıt için işaretleme
sub vcl_hit {
    # Cache'den gelen yanıt için işaret koy
    set req.http.X-Cache-Status = "HIT";
}

# Cache'de bulunamayan bir yanıt için işaretleme
sub vcl_miss {
    # Backend'den gelen yanıt için işaret koy
    set req.http.X-Cache-Status = "MISS";
}

# İstemciye yanıt gönderilmeden önce header ekle
sub vcl_deliver {
    # İşaretlenen header'ı istemciye gönder
    set resp.http.X-Cache = req.http.X-Cache-Status;
}