教材: [実践Docker ソフトウェアエンジニアの「Dockerよくわからない」を終わりにする本](https://zenn.dev/suzuki_hoge/books/2022-03-docker-practice-8ae36c33424b59)

---

## Dockerについての前提知識

- Docker DesktopはDocker EngineやLinuxカーネルが含まれているため、WindowsやMacでもDockerを利用できるようにしたもの。
  - WindowsかMacでDockerを使おうと思ったときに、インストールするDocker一式が入ったGUIアプリケーション
  - Docker Composeやk8sも入っている（知らなかった）
- ECSについては、「Docker Engineの入ったLinuxホストで、コンテナをそのままデプロイできる場所」ぐらいに捉えておく
- 「コンテナ≠仮想OS」であり **「コンテナ≒LinuxのNamespaceという機能によりほかと分離された、ただの1プロセス」**
  - Namespaceによって隔離されているため、ファイル名の衝突などが回避されている
  - Namespaceをイメージから作り出すことで異なるOSに見えるようにしてくれたり、Namespaceを簡単に作ったり消したりできるようなコマンドを提供してくれたりするものが、Dockerと言える
- イメージとは？
  - コンテナの実行に必要なパッケージで、ファイルやメタ情報を集めたもの
  - **複数のレイヤーというものからなる情報のことで、ホストマシンのどこかに `.img` のような具体的な単一のファイルが存在するわけではない**
  - イメージにはレイヤーによって次のような情報が含まれている
    - ベースはなにか？
    - 何をインストールしてあるか？
    - 環境変数はどうなっているか？
    - どういう設定ファイルを配置しているか？
    - デフォルト命令は何か？
- Dockerfileは **既存のイメージにレイヤーをさらに積み重ねる** ためのテキストファイル
  - インターネット上に公開されているイメージでは、インストールしてあるコマンドが足りないなどの問題がある場合に、自分に都合のいいイメージを作るためにDockerfileを作成する
- 基本のコマンドは以下の3つ！
  - **`container run`: イメージからコンテナを1つ作るだけのコマンド**
    - 以下の複数コマンドを一括で行っている
      - `image pull`
      - `container create`
      - `container start`
  - **`image build`: Dockerfileからイメージを作成するコマンド**
  - **`container exec`: コンテナに命令を送るコマンド**
    - `container stop` なども「コンテナをどうにかする」タイプに属するコマンド
- **イメージには、コンテナを起動したときにどんなコマンドを実行するかが予め書き込まれている**
  - ここではこれを「デフォルト命令」と呼ぶ
- **コンテナはある1つのコマンドを実行するために起動している**
  - **それはデフォルト命令か指定命令のどちらかで、そのPIDは1になる**
  - 複数のコンテナのPIDが1のものは、LinuxのNamespace機能により衝突しない
  - ここでは、PIDが1のプロセスをメインプロセス、それを起動するコマンドをメインコマンドと呼ぶ
  - **「コンテナはメインプロセスを実行するために起動する」**という点を押さえておくと、コンテナのことをよりスムーズに理解できるようになる。
    - 例えば、DBやWebサーバをひとまとめにしたコンテナを作ろう、という発想を避けられる
- コンテナが停止する理由は2つ
  - コンテナを停止する
  - **メインプロセスが終了する**
    - 例えば、メインプロセスを `ls` などにすると、コンテナ立ち上げ後、すぐに停止してしまう
- コンテナの状態変更を別のコンテナに反映するには？
  - 構成変更を全コンテナに反映したいなら、イメージを作る
  - コンテナのファイルを残したいならホストマシンと共有する
- `container exec` は、コンテナを活用したりトラブルシューティングに必要になったりするときにとても重要なコマンド
  - 具体的には以下のようなときに使いたい
    - コンテナの中にあるログを調べたい
    - **Dockerfileを書く前に `bash` でインストールコマンドを試し打ちしたい**
    - MySQLデータベースサーバのクライアント `mysql` を直接操作したい
- Dockerイメージの指定には、 `REPOSITORY:TAG` 形式を用いるのがわかりやすい
  - TAGを省略すると、 `latest` を指定したのとおなじになる
- Docker Hubのレイヤーページで確認できるものは、「Dockerfileによって積み上げられたレイヤーの情報」であり「Dockerfileそのもの」ではないことに注意！

---

## Dockerfileについて

Dockerfileはimageと紐付けてメンタルモデルを作っておくこと。

- Dockerfileの命令で、重要なもの
  - `FROM`: ベースイメージを指定する
  - `RUN`: 任意のコマンドを実行する
  - `COPY`: ホストマシンのファイルをイメージに追加する
  - `CMD`: デフォルト命令を指定する
- `RUN` で `apt` を使うのか `yum` を使うのかなどは、ベースイメージ（OS）によって決まる
  - `nginx:latest` など、一見してわからない場合は、bashを一度起動してOSやパッケージマネージャを調べることが、都度必要になる
- imageコマンド2つ
  - `build [option] <path>`: pathはDockerfileファイル内の `COPY` コマンドで使うファイルを指定するときの相対パス
    - **Dockerfileはディレクトリを分けてデフォルト名で管理するのが普通**
    - `-f, --file`: Dockerfileを指定する
      - 複数のDockerfileを使い分ける場合に使う
    - `-t, --tag`: ビルド結果にタグを付ける。人間が把握しやすいようにする
  - `history [option] <image>`: イメージはタグやイメージIDで指定

---

## container runとよく使われるオプション

- `-i, --interactive`: コンテナの標準入力に接続する
- `-t, --tty`: 疑似ターミナルを割り当てる
- `-d, --detach`: バックグラウンドで実行する
- `--rm`: 停止済みコンテナを自動で削除する（起動時に、停止済みコンテナのIDなどの重複を気にしなくて良くするために使う）
- `--name`: コンテナに名前をつける
  - コンテナ名は、停止済み・起動中のものを含めて、全体で一意である必要がある
- `--platform`: イメージのアーキテクチャを明示する（M1 Macで必要な場合がある、らしい）
  - WindowsおよびIntel Macは `linux/amd64` が、M1 Macの場合は `linux/arm64/v8` がDockerによって選ばれ、pullされる
    - 後者は、まだイメージが用意されていないものも多いらしい
    - その場合は、強引に `linux/amd64` を明示して選ぶ必要があるっぽい
- `-p, --publish <local>:<container>`: ポートの公開、これまでpはportの略だと思っていた

※ `--interactive, --tty` は `exec` でも使える

### コマンドヒストリー

```sh
# docker hubからイメージをpullしてからcontainer create & startまで行ってくれる
docker container run --publish 8080:80 nginx:1.21

# container stopの引数はIDかNAMEどちらでもよい
# container rmも同じ

# container ls --all で停止中のコンテナも表示
docker container ls --all

# rm --forceで停止と削除を一括で行う（コンテナ再起動は基本的にやらないので、よく使うことになりそう）

# ubuntuコンテナの対話シェルを呼ぶ
docker container run --interactive --tty ubuntu:20.04

# commandオプションを指定してrunすると、デフォルト命令を上書き？できる
docker container run --name nginx-bash --rm --interactive --tty nginx:1.21 bash

# execの例
docker container run --name ubuntu1 --rm --interactive --tty ubuntu:20.04
echo 'hello world' > ~/hello.txt
# 失敗してしまうような。。？
docker container exec \
         ubuntu1 \
         cat ~/hello.txt
# bashにつないで継続的にデバッグしたい場合
docker container exec --interactive --tty ubuntu1 bash

# イメージのビルド、タグがないと不便なので必ずつけること
docker image build --tag my-ubuntu:date .
docker container run --name my-ubuntu1 --rm my-ubuntu:date
docker container run --name my-ubuntu2 --rm --interactive --tty my-ubuntu:date vi
```

---

## 3部のハンズオンをやりながら

- `RUN` を書く場合に求められるのは、Dockerの知識ではなくLinuxの知識であることが大半！
  - 検索するときはイメージのOSやインストールなどのキーワードを使ったほうがよい
- `RUN` でなんでもかんでも `&&` でつなげてレイヤーを小さくするのが常に最適ではない、可読性を優先して分割することもある
- 構築手順が明らかではない状態からは、いきなりDockerfileは書かない
  - ベースイメージだけを起動して、そのあとbashで試行錯誤するのが普通
- Dockerfileを書くときは、、
  - **`FROM` を書くときは、Docker Hubで探す**
  - **`RUN` を書くときは、Linuxの知識が必要になる**
  - **`COPY` を書くときは、その製品の知識が必要になる**
  - 手順が明らかでない場合は、まずはコンテナを起動して手作業してみるのが有効（泥臭くやっていく！）
- Ubuntuなどの汎用イメージを除き、MySQLやNginxやRailsのような特定のサービスがセットアップされているイメージは、基本的にデフォルト命令で起動する
- 環境変数は `docker container run` で指定できる
- 一切のSQLを書かずにユーザやデータベースが作られたのは `mysql:5.7` のイメージに含まれるShell Scriptのおかげ。
  - 便利な機能の由来が（Dockerなのか特定のイメージなのか）どこかを意識することが大事！
- コンテナにbashがない場合はshを使う
  - 軽量であることを重視したAlpine Linuxをベースとしたイメージにはcurlやbashも入っていない場合が大半
- デバッグ時は `--detach` オプションを外して出力を隠さないようにすると良い
  - （個人的には、ログは細かくチェックしたいので、ターミナルを用意してdetachなしをデフォルトとしたほうが良さそう）

### Dockerのボリュームという概念

- ボリュームは、コンテナ内のファイルをホストマシン上でDockerが管理してくれる仕組み
  - **ホストマシン側のどこに保存されているかは関心がなく、とにかくデータを永続化したい** という場合に有用で、例えばデータベースのデータの永続化に活用できる
  - `docker volume create [option]`
  - `--name` オプションでの命名はほぼ必須
  - `container run` で作成済みボリュームをマウントするときのオプションは `--volume, --mount` の2つがあるが、記法がそれぞれ異なる
    - `--mount` のほうが可読性が高くてよい
  - `docker volume inspect docker-practice-db-volume` というコマンドで、volumeの詳細を調べられる
    - **MacなどのDocker Desktopを使っている場合、これはDocker DesktopのLinux上のパスなので、ホストマシンのパスとは全くの別物！**
      - **なので、ボリュームの削除はホストマシンには影響が一切ない**
- バインドマウントは、ホストマシンの任意のディレクトリをコンテナにマウントする仕組み
  - **ホストマシンとコンテナ双方がファイルの変更に関心がある、という場合に有用**
    - ソースコードの共有など
  - バインドマウントも、container run時には `--volume, --mount` オプションが両方使える
  - バインドマウントは、ホストマシンとファイルを共有するため、必要なとき以外は使わない
    - **まずはボリュームマウントを検討するべき！**
  - `COPY` との使い分け
    - `COPY` は `image build` をするタイミングでイメージにファイルを含めるので、コンテナが起動すればファイルが存在する
    - `COPY` の用途
      - 設定ファイルなど、コンテナによって変えない、かつ、めったに変更しないものを配置する場合
      - **本番デプロイ時のソースコードなど、即起動できる配布物を作る場合**
    - バインドマウントの用途
      - 開発時のソースコードなど、ホストマシンで変更したいがコンテナに随時反映したいものがある場合
      - 初期化クエリなど、イメージを配布する時点では用意できないものがある場合
    - **それぞれ、「イメージに対して行っている」か「コンテナに対して行っている」かが明確に違う** ということを意識すべし！

### Dockerのネットワーク

- **Dockerのコンテナはネットワークドライバというもので、Dockerネットワークに接続される**
  - ネットワークドライバには、ブリッジネットワークやオーバーレイネットワークがある
    - ブリッジネットワーク: Dockerのデフォルト、同一のDocker Engine上のコンテナが互いに通信をする場合に利用する
    - オーバレイネットワーク: 異なるDocker Engine上のコンテナが互いに通信をする場合に利用する
- **デフォルトブリッジネットワーク**
  - コンテナを起動する際に、ネットワークドライバについて一切の指定を行わないと、デフォルトブリッジネットワークが自動的に生成され、コンテンはこのネットワークに接続される
  - 特徴
    - コンテンが通信するためには、すべてのコンテナ間をリンクする操作が必要になる
    - コンテナ間の通信はIPアドレスで行う
    - Docker Engine上のすべてのコンテナ（例えば別プロジェクト）に接続できてしまう
- **ユーザ定義ブリッジネットワーク**
  - 特徴
    - **相互通信をできるようにするには、同じネットワークを割り当てるだけで良い**
    - **コンテナ間で自動的にDNS解決を行える**
      - `--network-alias` オプションで設定する名前を使う？
    - 通信できるコンテナが同一ネットワーク上のコンテナに限られ、隔離度が上がる
  - （このあたり、すごい疑問に思ってたことが解説されていた！）
  - 具体的には、containerおよびnetworkのinspectコマンドを見ると、ゲートウェイなどを確認できる
- 「コンパイルをしてほしい」とか「静的コンテンツをホスティングしてほしい」のような、コンテナ間通信を必要としない場合は、デフォルトブリッジネットワークで十分

### コマンドリスト

```sh
# appコンテナのイメージ確認
docker container run --name app --rm --interactive --tty docker-practice:app bash
docker container run --name app --rm docker-practice:app php -S 0.0.0.0:8000
docker container exec --interactive --tty app bash

# dbコンテナの起動
docker container run \
         --name db \
         --rm \
         --platform linux/amd64 \
         --env MYSQL_ROOT_PASSWORD=master-pass \
         --env MYSQL_USER=main \
         --env MYSQL_PASSWORD=main-pass \
         --env MYSQL_DATABASE=event \
         docker-practice:db
docker container exec --interactive --tty db bash
mysql -h localhost -u main -pmain-pass event

# mailコンテナの起動
docker container run \
         --name mail \
         --rm \
         --platform linux/amd64 \
         mailhog/mailhog:v1.0.1
docker container exec \
         --interactive \
         --tty \
         --user root \
         mail \
         sh

# volumeの作成とdbコンテナへのマウント
docker volume create \
         --name docker-practice-db-volume
docker container run                                                     \
    --name db                                                            \
    --rm                                                                 \
    --detach                                                             \
    --platform linux/amd64                                               \
    --env MYSQL_ROOT_PASSWORD=rootpassword                               \
    --env MYSQL_USER=hoge                                                \
    --env MYSQL_PASSWORD=password                                        \
    --env MYSQL_DATABASE=event                                           \
    --mount type=volume,src=docker-practice-db-volume,dst=/var/lib/mysql \
    docker-practice:db

# バインドマウントしつつcontainer run、typeはbindとする
docker container run                          \
    --name app                                \
    --rm                                      \
    --detach                                  \
    --interactive                             \
    --tty                                     \
    --mount type=bind,src=$(pwd)/src,dst=/src \
    docker-practice:app                       \
    php -S 0.0.0.0:8000 -t /src

# dbコンテナに初期化用クエリファイルをバインドマウントする
docker container run                                                                         \
    --name db                                                                                \
    --rm                                                                                     \
    --detach                                                                                 \
    --platform linux/amd64                                                                   \
    --env MYSQL_ROOT_PASSWORD=rootpassword                                                   \
    --env MYSQL_USER=hoge                                                                    \
    --env MYSQL_PASSWORD=password                                                            \
    --env MYSQL_DATABASE=event                                                               \
    --mount type=volume,src=docker-practice-db-volume,dst=/var/lib/mysql                     \
    --mount type=bind,src=$(pwd)/docker/db/init.sql,dst=/docker-entrypoint-initdb.d/init.sql \
    docker-practice:db

docker container exec --interactive --tty db mysql -h localhost -u hoge -ppassword event

# ポートの公開
docker container run                          \
    --name app                                \
    --rm                                      \
    --detach                                  \
    --interactive                             \
    --tty                                     \
    --mount type=bind,src=$(pwd)/src,dst=/src \
    --publish 18000:8000                      \
    docker-practice:app                       \
    php -S 0.0.0.0:8000 -t /src
docker container run       \
    --name mail            \
    --rm                   \
    --detach               \
    --platform linux/amd64 \
    --publish 18025:8025   \
    mailhog/mailhog:v1.0.1

# ユーザ定義ブリッジネットワークを作成する
docker network create \
         docker-practice-network
# ネットワークを利用する形でcontainer run
# 接続される側は --network-alias オプションでのエイリアス名の指定が必要？
docker container run                          \
    --name app                                \
    --rm                                      \
    --detach                                  \
    --mount type=bind,src=$(pwd)/src,dst=/src \
    --publish 18000:8000                      \
    --network docker-practice-network         \
    docker-practice:app                       \
    php -S 0.0.0.0:8000 -t /src
docker container run                                                                         \
    --name db                                                                                \
    --rm                                                                                     \
    --detach                                                                                 \
    --platform linux/amd64                                                                   \
    --env MYSQL_ROOT_PASSWORD=rootpassword                                                   \
    --env MYSQL_USER=hoge                                                                    \
    --env MYSQL_PASSWORD=password                                                            \
    --env MYSQL_DATABASE=event                                                               \
    --mount type=volume,src=docker-practice-db-volume,dst=/var/lib/mysql                     \
    --mount type=bind,src=$(pwd)/docker/db/init.sql,dst=/docker-entrypoint-initdb.d/init.sql \
    --network docker-practice-network                                                        \
    --network-alias db                                                                       \
    docker-practice:db
docker container run                  \
    --name mail                       \
    --rm                              \
    --detach                          \
    --platform linux/amd64            \
    --publish 18025:8025              \
    --network docker-practice-network \
    --network-alias mail              \
    mailhog/mailhog:v1.0.1
# pingで疎通確認
docker container exec --interactive --tty app ping db -c 3
```

---

## Docker Compose

- Docker Composeは `image build` と `container run` の手順を置き換える
  - Dockerfileは依然として必要
    - ※Docker Hubのイメージをそのまま使う場合は不要
  - コンテナの起動を楽にするツールと考える
- `container_name` は、無理に短くせず、わかりやすい命名を行う
- イメージのビルドは、 `compose up` をしたがイメージが存在しない場合、もしくは `compose up` に `build` オプションをつけて実行した場合に行われる
- **Docker Composedeは、自動でブリッジネットワークが作成される**
  - **また、サービス名が自動でコンテナのネットワークにおけるエイリアスとして設定される**
- `docker compose down` は、 `container run` の `--rm` オプションと同様、停止済みコンテナを削除する形で停止処理を行う

---

## デバッグノウハウ

- 覚えておきたいサブコマンド `ls, inspect`
  - image, container, volume, network全てに対して使えるサブコマンド
- `docker container logs` は、バックグラウンド起動したコンテナの出力を確認できる
  - `--follow` とすると `tail -f` のような感覚で使える
- Docker DesktopのGUIを見るのも良い
- サービスそのもののログを見るのも大事
  - サービスの設定ファイルなどから、どこにログが吐かれるのかは知っておきたい

