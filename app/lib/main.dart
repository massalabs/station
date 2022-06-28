import 'package:flutter/material.dart';
import 'package:webview_flutter_plus/webview_flutter_plus.dart';
import 'package:thyra_bridge_dart/thyra_bridge_dart.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const MyHomePage(),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key}) : super(key: key);
  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  
  @override
  Widget build(BuildContext context) {
    // String webContent = ''' <!DOCTYPE html>
    // <html>
    // <body>

    // <h2>Example to load html from string</h2>

    // <p>This is paragraph 1</p>

    // <img src="https://thumbs.dreamstime.com/b/sun-rays-mountain-landscape-5721010.jpg"  width="250" height="250">

    // </body>
    // </html>''';
    final binding = Binding();
    // String webContent = String.fromCharCodes(binding.fetchWoc("1234", "index.html"));
    String webContent = String.fromCharCodes(binding.fetchWebsite("A1aMywGBgBywiL6WcbKR4ugxoBtdP9P3waBVi5e713uvj7F1DJL", "index.html"));
    return Scaffold(
      appBar: AppBar(
        title: const Text('WebView Load HTML Locally From assets'),
      ),
      body: WebViewPlus(
        javascriptMode: JavascriptMode.unrestricted,
        onWebViewCreated: (controller){
          controller.loadString(webContent);
        },
      ),
    );
  }
}
