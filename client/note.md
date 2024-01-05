Reactアプリの新規作成時のコマンド
下記のコマンドではproxy環境下にないにも関わらずproxy関連のエラーにより失敗。

```
npx create-react-app my-app
cd my-app
npm start
```

下記のコマンドでは何故か成功した。

```
<!-- React -->
yarn global add create-react-app
<!-- 関数型ライブラリ -->
yarn add fp-ts
<!-- グローバルインストールされたツールが利用できない場合は下記コマンドを実行する -->
export PATH="$PATH:`yarn global bin`"
<!-- Reactアプリの作成 -->
yarn create react-app app
cd app
<!-- Reactアプリの実行 -->
yarn start
```

```
emotion
fp-ts
react
typescript
```