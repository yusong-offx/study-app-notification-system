import 'package:web_socket_channel/web_socket_channel.dart';

Future<void> websocketCummunication(Function callback) async {
  final wsUrl = Uri.parse('ws://127.0.0.1:3000/ws/yusong-offx');
  final ws = WebSocketChannel.connect(wsUrl);

  ws.stream.listen((msg) {
    callback(msg);
  });

  ws.sink.add('connect');
  print('websocket connect!');
}
