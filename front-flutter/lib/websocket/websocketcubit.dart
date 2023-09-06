import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:study_notification/websocket/websocket.dart';

class WebSocketMessageCubit extends Cubit<String> {
  WebSocketMessageCubit() : super('');

  void connect() {
    websocketCummunication((msg) {
      emit(msg);
    });
  }

  void sendMessage(String message) {
    emit(message);
  }

  void websocketsendMessage(String message) {
    emit(message);
  }
}
