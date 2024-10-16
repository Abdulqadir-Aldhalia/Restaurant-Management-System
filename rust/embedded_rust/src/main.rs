#![no_std]
#![no_main]

// use core::result::Result;
// use embedded_svc::mqtt::client::{Event, QoS, Connection};
// use embedded_svc::wifi::{AuthMethod, ClientConfiguration, Configuration};
// use esp_idf_hal::peripherals::Peripherals;
// use esp_idf_svc::eventloop::EspSystemEventLoop;
// use esp_idf_svc::mqtt::client::{EspMqttClient, MqttClientConfiguration};
// use esp_idf_svc::nvs::EspDefaultNvsPartition;
// use esp_idf_svc::wifi::{BlockingWifi, EspWifi};
// use esp_idf_sys::esp_println as println;
// use esp_idf_sys::link_patches;

#[no_mangle]
fn main () {

}
// fn main() -> anyhow::Result<()> {
//     link_patches();

//     // Initialize peripherals
//     let peripherals = Peripherals::take().unwrap();
//     let sysloop = EspSystemEventLoop::take()?;
//     let nvs = EspDefaultNvsPartition::take()?;

//     // Initialize Wi-Fi
//     let mut wifi = BlockingWifi::wrap(
//         EspWifi::new(peripherals.modem, sysloop.clone(), Some(nvs))?,
//         sysloop,
//     )?;

//     wifi.set_configuration(&Configuration::Client(ClientConfiguration {
//         ssid: "Wokwi",             // Set your SSID here
//         password: "password123",   // Set your Wi-Fi password here
//         auth_method: AuthMethod::WPA2Personal,
//         ..Default::default()
//     }))?;

//     // Wait for the Wi-Fi connection
//     wifi.start()?;
//     if wifi.is_connected().unwrap() {
//         let config = wifi.get_configuration().unwrap();
//         println!("Connected to Wi-Fi: {:?}", config);
//     } else {
//         println!("Failed to connect to Wi-Fi");
//         return Err(anyhow::anyhow!("Wi-Fi connection failed"));
//     }

//     println!("Wi-Fi Connected");

//     // MQTT client configuration
//     let mqtt_config = MqttClientConfiguration::default();

//     // Initialize the MQTT client
//     let mut client = EspMqttClient::new(
//         "tcp://broker.emqx.io:1883", // Set your MQTT broker URL
//         &mqtt_config,
//         move |message_event| {
//             match message_event {
//                 Ok(Event::Connected(_)) => println!("Connected to MQTT broker"),
//                 Ok(Event::Subscribed(id)) => println!("Subscribed with ID: {}", id),
//                 Ok(Event::Received(msg)) => {
//                     if !msg.data().is_empty() {
//                         println!("Received: {:?}", core::str::from_utf8(msg.data()).unwrap());
//                     }
//                 }
//                 _ => println!("Unhandled MQTT event: {:?}", message_event),
//             };
//         },
//     )?;

//     // Subscribe to an MQTT topic
//     client.subscribe("esp32/notify", QoS::AtLeastOnce)?;

//     // Keep the system alive and prevent watchdog timeout
//     loop {
//         esp_idf_hal::delay::FreeRtos::delay_ms(1000);
//     }
// }
