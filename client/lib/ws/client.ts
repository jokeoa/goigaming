import { WS_URL } from "@/lib/constants";

type WSClientOptions = {
  readonly url?: string;
  readonly token?: string;
  readonly gameId?: string;
  readonly onMessage?: (data: unknown) => void;
  readonly onOpen?: () => void;
  readonly onClose?: () => void;
  readonly onError?: (error: Event) => void;
  readonly reconnect?: boolean;
  readonly reconnectInterval?: number;
  readonly maxReconnectAttempts?: number;
};

export class WSClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private readonly options: Required<
    Pick<
      WSClientOptions,
      "reconnect" | "reconnectInterval" | "maxReconnectAttempts"
    >
  > &
    WSClientOptions;

  constructor(options: WSClientOptions) {
    this.options = {
      reconnect: true,
      reconnectInterval: 3000,
      maxReconnectAttempts: 5,
      ...options,
    };
  }

  connect(): void {
    const params = new URLSearchParams();
    if (this.options.token) {
      params.set("token", this.options.token);
    }
    if (this.options.gameId) {
      params.set("game_id", this.options.gameId);
    }

    const qs = params.toString();
    const baseUrl = this.options.url ?? WS_URL;
    const url = qs ? `${baseUrl}?${qs}` : baseUrl;

    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.options.onOpen?.();
    };

    this.ws.onclose = () => {
      this.options.onClose?.();
      this.attemptReconnect();
    };

    this.ws.onerror = (event) => {
      this.options.onError?.(event);
    };

    this.ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this.options.onMessage?.(data);
      } catch {
        // ignore malformed messages
      }
    };
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  send(data: unknown): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data));
    }
  }

  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  private attemptReconnect(): void {
    if (
      !this.options.reconnect ||
      this.reconnectAttempts >= this.options.maxReconnectAttempts
    ) {
      return;
    }

    this.reconnectAttempts += 1;
    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, this.options.reconnectInterval);
  }
}
