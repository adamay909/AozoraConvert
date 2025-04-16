# AozoraConvert

Provides tools for converting texts from Aozora Bunko to epub and kindle documents formatted for vertical typesetting. For more documentation in English, check out the Go documentations, in particular the one for cmd/azrconvert.

## 概要

青空文庫で提供されているテキストを縦書きのWebページ、あるいは縦書きのEpub3及びKindle用のAZW3（KF8)に変換するためのGoライブラリとコマンドラインツール：
- 変換後のエンコーディングはすべてUTF8。
- 外字も可能な限りUTF8に変換し、画像を避ける。
- 挿絵なども自動的に取り込む。
- EpubとKindle用では簡単な表紙を追加。

Kindle用について：
- kindlegenを必要としない。
- Kindleストアから入手できるものと違いDRMは一切なし。

 

## コマンドラインツール

cmd/azrconvertフォルダにコンパイル済みのバイナリがいくつか置いてある。もちろん同フォルダで go build で自分でバイナリ作成も可。
- azrconvert （Linux/amd64用）
- azrconvert.exe （Windows用）
- azrconvert_APPLESILICON （M1以降のMac用）
- azrconvert_INTEL （Intelチップ仕様のMac用）

Linux/amd64以外はクロスコンパイルしたもの。

### 使い方

芥川龍之介の「芋粥」をEpub3に変換する： 
```
$ azrconvert -epub https://www.aozora.gr.jp/cards/000879/files/55_14824.html
```
  -epub オプションによってEpubに変換する。URLは青空文庫の「いますぐXHTML版で読む」をクリックしたときに表示されるページのURL。出力は
 ```
Output written to 芋粥.epub.
  ```
  
 芋粥.epubというファイルが作成されたことがわかる。
  
オプションとして -kindle を指定すればKindle用のAZW3ファイルが作成される。

-web を指定すると、縦書き用のHTMLファイルとCSS、及び画像ファイルをパッケージしたZIPアーカイブが作成される。アーカイブを解凍して、その中の1.htmlを縦書き表示対応のブラウザで開けばテキストが縦書きで表示される（最近のFirefox, Google Chrome, Safariはいずれも問題なし）。

EpubとKindle用を同時に作成することもできる：
```
$ azrconvert -epub -kindle  https://www.aozora.gr.jp/cards/000879/files/55_14824.html

Output written to 芋粥.epub.
Output written to 芋粥.azw3.
```

無論 -web も一緒に指定することができる。

出力されるファイル名は\<meta name="DC.Title">によってタイトルが指定されている場合はそれを利用する。それができない場合はURLのベースを利用する。例えば田中正造の「公益に有害の鑛業を停止せざる儀に付質問書」の場合は次のようになる：
```
$ azrconvert -kindle https://www.aozora.gr.jp/cards/000649/files/4958_10239.html

Output written to 4958_10239.azw3.
```

手動で指定したい場合は-oオプションを使う：
```
$ azrconvert -web -o imogayu  https://www.aozora.gr.jp/cards/000879/files/55_14824.html

Output written to imogayu.zip.
```
いずれにせよコマンドプロンプトへの出力で出力ファイル名を確認できる。

-v オプションを使うとlogを画面とazrconvert.logの双方に出力する。基本的に必要ない。

## 留意点

- 最近のブラウザ（Firefox, Google Chrome, Safari)はいずれも縦書きの日本語ページを問題なく表示できるが、使用するフォントによってはうまく行かないので、表示がおかしかったらまずフォントをかえてみること。Noto Serif JP、Noto Sans JP、 IPAフォントなどは大丈夫。
- 使用しているCSSはFont Familyをserif, sans-serifの順で指定しているので、表示フォントをかえるにはブラウザの設定でserifの方をかえる。
- 大概の場合、使用に耐えるものを作成できるが、Epub等の微調整をしたい場合は万能ツールの[Calibre](https://calibre-ebook.com/ja/download)の使用がお薦め。
- 電子ブックリーダーにファイルを送るのも[Calibre](https://calibre-ebook.com/ja/download)がお薦め。


