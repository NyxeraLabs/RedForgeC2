pub trait Transport {
    /// Sends a raw payload to the C2 channel.
    fn send(&mut self, payload: &[u8]) -> std::io::Result<()>;

    /// Receives a raw payload from the C2 channel.
    fn receive(&mut self) -> std::io::Result<Vec<u8>>;

    /// Closes the transport and cleans up resources.
    fn close(&mut self) -> std::io::Result<()>;
}
