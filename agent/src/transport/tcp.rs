use crate::transport::Transport;
use std::io::{Read, Write};
use std::net::{TcpStream, ToSocketAddrs};
use std::time::Duration;

/// A simple TCP transport implementation.
///
/// This is primarily here as a skeleton for future transport layers (e.g. raw
/// socket communication) and as an example of a non-HTTP transport option.
pub struct TcpTransport {
    stream: TcpStream,
}

impl TcpTransport {
    /// Connects to the given address (e.g. "example.com:1234").
    pub fn connect<A: ToSocketAddrs>(addr: A, timeout: Duration) -> std::io::Result<Self> {
        let stream = TcpStream::connect_timeout(&addr.to_socket_addrs()?.next().ok_or_else(|| {
            std::io::Error::new(std::io::ErrorKind::InvalidInput, "invalid socket address")
        })?, timeout)?;
        stream.set_read_timeout(Some(timeout))?;
        stream.set_write_timeout(Some(timeout))?;
        Ok(Self { stream })
    }
}

impl Transport for TcpTransport {
    fn send(&mut self, payload: &[u8]) -> std::io::Result<()> {
        // Send the length-prefix and payload.
        let len = (payload.len() as u32).to_be_bytes();
        self.stream.write_all(&len)?;
        self.stream.write_all(payload)?;
        self.stream.flush()?;
        Ok(())
    }

    fn receive(&mut self) -> std::io::Result<Vec<u8>> {
        let mut len_buf = [0u8; 4];
        self.stream.read_exact(&mut len_buf)?;
        let len = u32::from_be_bytes(len_buf) as usize;
        let mut buf = vec![0u8; len];
        self.stream.read_exact(&mut buf)?;
        Ok(buf)
    }

    fn close(&mut self) -> std::io::Result<()> {
        self.stream.shutdown(std::net::Shutdown::Both)
    }
}
