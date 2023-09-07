import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:study_notification/websocket/websocket.dart';

class WebSocketMessageCubit extends Cubit<String> {
  WebSocketMessageCubit() : super('');

  void connect() {
    print("websocket connect!");
    websocketCummunication((msg) {
      emit(msg);
    });
  }

  void sendMessage(String message) {
    emit(message);
  }
}
