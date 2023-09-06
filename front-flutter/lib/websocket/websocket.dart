import 'package:web_socket_channel/web_socket_channel.dart';

Future<void> websocketCummunication(Function callback) async {
  final wsUrl = Uri.parse('ws://localhost:8080/v1');
  final ws = WebSocketChannel.connect(wsUrl);

  ws.stream.listen((msg) {
    ws.sink.add('Hello from Flutter');
    callback(msg);
  });
}
