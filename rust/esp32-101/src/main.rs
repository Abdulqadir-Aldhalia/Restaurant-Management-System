#![no_std]
#![no_main]

use esp_backtrace as _;
use esp_hal::{
    clock::ClockControl,
    delay::Delay,
    peripherals::Peripherals,
    prelude::*,
    gpio::{Input, PullUp, Gpio12},
    system::SystemControl,
};
use panic_halt as _; // panic handler that halts the program
use embedded_hal::digital::v2::InputPin;
use smoltcp::iface::{EthernetInterfaceBuilder, NeighborCache};
use smoltcp::wire::{EthernetAddress, IpCidr, Ipv4Address};
use ureq; // For HTTP requests
use esp_wifi::{WifiController, initialize, WifiState, ClientConfiguration};

// Wi-Fi and network configurations (stub)
const SSID: &str = "mySSID";
const PASSWORD: &str = "myPASSWORD";

#[entry]
fn main() -> ! {
    // Initialize peripherals
    let peripherals = Peripherals::take().unwrap();

    // Set up the system clocks
    let mut system = peripherals.SYSTEM.split();
    let clocks = ClockControl::boot_defaults(system.clock_control).freeze();

    // Create a delay instance for timed operations
    let mut delay = Delay::new(&clocks);

    // GPIO configuration: Button connected to GPIO12
    let gpio = peripherals.GPIO.split();
    let button: Gpio12<Input<PullUp>> = gpio.gpio12.into_pull_up_input();

    // Initialize Wi-Fi and connect
    connect_to_wifi();

    // Infinite loop to check button press
    loop {
        // Check if the button is pressed
        if button.is_low().unwrap() {
            // Debounce the button by adding a delay
            delay.delay_ms(50u32);

            // Action when the button is pressed
            println!("Button pressed! Sending notification to server...");

            // Send the HTTP request to notify the server
            send_notification_to_server();

            // Add a delay to avoid multiple triggers
            delay.delay_ms(1000u32);
        }
    }
}

// Wi-Fi connection function using esp-wifi
fn connect_to_wifi() {
    println!("Initializing Wi-Fi...");

    // Initialize Wi-Fi with default settings
    let (wifi, _modem) = initialize().unwrap();
    
    // Create a Wi-Fi controller
    let mut wifi_controller = WifiController::new(wifi);

    // Set up Wi-Fi client configuration
    let client_config = ClientConfiguration {
        ssid: SSID.into(),
        password: PASSWORD.into(),
        ..Default::default()
    };

    // Apply the configuration and connect
    wifi_controller.set_configuration(&client_config).unwrap();
    wifi_controller.start().unwrap();

    // Wait for the Wi-Fi to connect
    while wifi_controller.get_state() != WifiState::Connected {
        println!("Waiting for Wi-Fi connection...");
        // Add a small delay to avoid busy-waiting
        delay.delay_ms(500u32);
    }

    println!("Connected to Wi-Fi!");
}

// Function to send an HTTP POST request using ureq
fn send_notification_to_server() {
    // Define the API key and server endpoint
    const API_KEY: &str = "[Q<-(C*V{u/AJim+<qwJ0|~Jus{u',pYJ]vEflDl~sb5LiLx2JA}F,.&cJB'a{u";
    let url = "http://localhost:8000/embeddedSystem/notify?table_id=5f7d6cbb-46e5-4fe8-9cfa-24ceee2a18a9";

    // Sending the HTTP POST request
    let response = ureq::post(url)
        .set("X-API-Key", API_KEY)
        .call();

    match response {
        Ok(res) => {
            // Print the response status
            println!("Server response: {}", res.status());
        }
        Err(e) => {
            // Handle error
            println!("Failed to send request: {:?}", e);
        }
    }
}

// Minimal setup for smoltcp interface (you'd need to configure a real interface)
// Example of what would be involved in setting up smoltcp TCP networking stack:
fn setup_networking() -> EthernetInterfaceBuilder<'static, 'static> {
    let ip_addrs = [IpCidr::new(Ipv4Address::new(192, 168, 1, 100).into(), 24)];
    let neighbor_cache = NeighborCache::new();
    let ethernet_addr = EthernetAddress([0x02, 0x00, 0x00, 0x00, 0x00, 0x01]);

    EthernetInterfaceBuilder::new() // Assuming you have a network device available
        .ethernet_addr(ethernet_addr)
        .neighbor_cache(neighbor_cache)
        .ip_addrs(ip_addrs)
}
