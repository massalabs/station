import 'package:flutter/material.dart';
import 'package:webview_flutter_plus/webview_flutter_plus.dart';
import 'package:thyra_bridge_dart/thyra_bridge_dart.dart';

void main() {
  runApp(const MyApp());
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
  // String webContent = String.fromCharCodes(binding.fetchWebsite(
  //     "A1aMywGBgBywiL6WcbKR4ugxoBtdP9P3waBVi5e713uvj7F1DJL", "index.html"));

  late String webContent;
  late WebViewPlusController controller;
  late WebViewPlus body;

  @override
  void initState() {
    super.initState();
    webContent = "testttt";
    body = WebViewPlus(
        javascriptMode: JavascriptMode.unrestricted,
        onWebViewCreated: (WebViewPlusController webViewPlusController) {
          controller = webViewPlusController;
        }
        // onWebViewCreated: (controller) {
        //   controller.loadString("testttt");
        // },
        );
  }

  @override
  Widget build(BuildContext context) {
    // WebViewPlusController controller;

    return Scaffold(
      appBar: AppBar(
        // title: const Text('WebView Load HTML Locally From assets'),
        title: TextField(
          textAlign: TextAlign.center,
          // The style of the input field
          decoration: const InputDecoration(
            hintText: 'Title',
            icon: Icon(Icons.edit), // Edit icon
            // The style of the hint text
            hintStyle: TextStyle(
              color: Colors.black,
              fontSize: 18,
            ),
          ),
          onSubmitted: (text) {
            _doSomething(text, controller);
          },
          // The controller of the input box
        ),
      ),
      // body: _getWebView(webContent),
      body: body,
    );
  }

  _getWebView(webContent) {
    print('lololol');
    print(webContent);
    return WebViewPlus(
      javascriptMode: JavascriptMode.unrestricted,
      // onWebViewCreated: (WebViewPlusController webViewPlusController) {
      //   controller = webViewPlusController;
      // }
      onWebViewCreated: (controller) {
        controller.loadString(webContent);
      },
    );
  }

  void _doSomething(String text, WebViewPlusController controller) {
    // Using the callback State.setState() is the only way to get the build
    // method to rerun with the updated state value.
    final binding = Binding();
    setState(() {
      if (text == "lol") {
        webContent = String.fromCharCodes(binding.fetchWebsite(
            "A1MrqLgWq5XXDpTBH6fzXHUg7E8M5U2fYDAF3E1xnUSzyZuKpMh",
            "index.html"));
        print('lol');
        print(webContent);
        controller.loadString(webContent);
      } else {
        webContent = String.fromCharCodes(binding.fetchWebsite(
            "A1aMywGBgBywiL6WcbKR4ugxoBtdP9P3waBVi5e713uvj7F1DJL",
            "index.html"));
        print('else');
        print(webContent);
        controller.loadString(webContent);
      }
    });
  }
}
