# Apache Kafka

#### O kafka além de uma ferramenta é um ecosistema

## Producer: Garantia de entrega
```
		   - Broker A -> leader
- Producer - Broker B -> follower
		   - Broker C -> follower

```

### Sempre entregue ao líder
	- Parametros:
		- Ack 0;
			- 0: none, ff (fire and forget), não precisa me informar que gravou a mensagem;
			- Caso de geolocalização do uber por exemplo;
		- Ack 1:
			- 1: leader, informa o producer que salvou a mensagem;
			- Mais lento, garantia de entrega ainda que não em todos os casos;
		- Ack -1:
			- "-1": leader salva, replica, os followers salvam, avisam o leader e ele retorna dizendo que tá tudo suave cpx ta guardado por Deus;
			- Muito mais lento;

### Producer: Garantia de entrega
	- At most once: Melhor performance. Pode perder algumas mensagens;
		- [1,2,3,4,5] -> Kafka process -> [1,3,4,5];
	- At least once: Performance moderada. Pode duplicar mensagens
		- [1,2,3,4,5] -> Kafka process -> [1,2,3,4,4,5];
	- Exacly once: Pior performance. Exatamente uma vez:
		- [1,2,3,4,5] -> Kafka process -> [1,2,3,4,5];

### Producer: Indepotência
	- OFF: quando ocorre um erro de rede por exemplo pode ocorrer a duplicação de mensagem em detrimento do duplo envio por parte do produtor;
	- ON: Kafka consegue descartar a duplicada e inserir a mensagem na order correta por ter os timestamps das mensagens;

## Consumer: 

- Producer -> Topic (pt0,1,2) -> consumer
	- O consumer lê a pt 0,1 e 2 pois só tenho 1 consumer;

- Consumer groups:
	- Producer -> Topic (pt0,1,2) -> group x (consumerA, consumerB);
		- (consumerA, consumerB) podem ser o mesmo software rodando em máquinas diferentes;
		- consumerA lê pt0,1
		- consumerB lê pt2
	- Se tivesse um consumer fora do grupo ele leria 0,1,2, só ocorre essa distribuição quando há grupo;
	- Producer -> Topic (pt0,1,2) -> group x (consumerA, consumerB, consumerC);
		- consumerA lê pt0
		- consumerB lê pt1
		- consumerC lê pt2
	- Se tivesse um D ele ficária em iddle, só um consumer por partição

- #### Ideal: a mesma quantidade de partitions e consumers.

- #### Trabalhar com consumer groups é muito comum, um consumer sozinho é o grupo dele mesmo;

- Apache Kafka
	- Projeto que faz parte da apache foundation;
	- Plataforma distribuída de streaming de eventos open-sorce;

- O mundo dos eventos
	- Demanda por eventos crescente em função de IOT, monitoramento de aplicações, sistemas de alarme, etc.

- Perguntas acerca:
	- Onde eu posso salvar esses evnetos?
	- Como recuperar de forma rápida e simples de forma que o feedback entre sistemas aconteça de forma fluída em tempo real?
	- Como escalar?
	- Como consigo ter resiliência e alta disponibilidade?

- Kafka e seus super poderes
	- Altíssimo thoughput (capacidade de receber e processar requisições);
	- Latência extremamente baixa (2ms);
	- Escalável;
	- Armazenamento (banco de dados otimizado);
	- Se conecta com quase tudo;
	- Bibliotecas prontas para as mais diversas tecnologias;
	- Ferramentas open-source;

- Empresas usando Kafka
	- Linkedin;
	- Netflix;
	- Uber;
	- Twitter;
	- Dropbox;
	- Spotify;
	- Paypal;
	- Bancos...

- Dinâmica básica de funcionamento:

- Producer -> Kafka -> Broker <- Consumer 
	- Kafka é um cluster, um conjunto de máquinas, formado de nós
	- Cada Broker possui seu próprio DB;
	- O Kafka não envia nenhuma mensagem para ninguém, o consumer que quem vai atrás da mensagem;
	- Recomendação mínima é de 3 Brokers para subir um apache kafka em produção;

- Apache Zookeeper é usado em conjunto para o serviço de service discovery, porém, está sendo descontinuado para independência do Kafka;

- Topics:
	- É o canal de comunicação responsavél por receber e disponibilizar os dados enviados para o Kafka;
	- Kafka !== RabbitMQ;
	- As filas do RMQ são parecidas com os tópicos do Kafka;
	- Producer -> Kafka <- (consumer1, consumer2)
						

- Posso ter diversos sistemas diferentes lendo os mesmos tópicos e eles podem consumir as mesmas mensagem diferente do RMQ.

### Topic ~= log:
	
- Cada mensagem enviada ganha um id que é chamado de offset;
- As mensagens continuam disponíveis independente do tempo de leitura do sistema, sistema A pode ler hoje e sistema B amanhã;
	

- Anatomia de um registro:
```					
				 - Headers
				 - Key	
	- Offset 0 -> - Value
				 - Timestamp
```

- Partições:
	- Cada tópico pode ter uma ou mais partições para garantir a distribuição e resiliência de seus dados
	- Topic 
		- partition 1 (broker A)
		- partition 2 (broker B)
		- partition 3 (broker C)

- Economiza velocidade porque aloca máquinas para partições evitando desperdicio de leitura


- Sobre as "keys"
	- Sistema bancário:
		- Partition 1 -> [0,1 (transferencia)] -> consumer 1 (lento)
		- Partition 2 -> [0 (estorno)] -> consumer 2 (rápido)
		- Partition 3 -> [0] -> consumer 2

	- Garantindo a ordem seria: 
		- [0(transferencia),1 (estorno)] 
			- Garantir a entrega de mensagem entregando na mesma partição
	- Keys:
		- [0(transferencia, key=movimentacao),1 (estorno, key=movimentacao)]
			- Todas as mensagens com a mesma key vão para mesma partição, logo:
				- Quando preciso garantir a ordem eu utilizo uma key.


- Partições distribuídas:
```
				-> Broker A (partition 1)
- topic vendas  -> Broker B (partition 2)
			    -> Broker C (partition 3)

- Replication factor = 2

				-> Broker A (partition 1,3)
- topic vendas  -> Broker B (partition 2,1)
			    -> Broker C (partition 3,2)
```
- Usado para evitar problemas caso um Broker caia garantindo assim a risiliência;
	- De acordo com nível de criticidade da aplicação nós temos mais replicações;


- Partition Leadership
	- Existem 2 tipos de partition, as leaders e as followers;
	- O Consumer vai necessáriamente ler sempre a partição Leader;
	- Porque ter as partições então? 
		- Quando um Broker cai, ele verifica as partições que o mesmo tinha, verifica onde existe a cópia em outro broker e aí sim ele vai consumir deste broker transformando-a em uma leadership;
		- A leader sempre é a lida, não importa a situação.

### Criar topico

#### Listar topicos
``
kafka-topics --list --bootstrap-server=localhost:9092
``

#### Descrever topico
`` 
kafka-topics --topics --bootstrap-server=localhost:9092 --topic=teste --describe
`` 

#### Consumir mensagens 
``
kafka-console-consumer --bootstrap-server=localhost:9092 --topic=teste --from-beginning
`` 

#### Consumir mensagens com grupos 
`` 
kafka-console-consumer --bootstrap-server=localhost:9092 --topic=teste --group=x
`` 

#### Enviar mensagens
`` 
kafka-console-producer --bootstrap-server=localhost:9092 --topic=teste
`` 


### Descrever grupo
- kafka-consumer-groups --bootstrap-server=localhost:9092 --group=x --describe
	- GROUP  =  nome do grupo (dã)         
	- TOPIC  =  nome do tópico (dã)
	- PARTITION  = partição 
	- CURRENT-OFFSET  = apontamento atual, por exemplo, tem 30 mensagens, tu leu 10 então o apontamento é 10
	- LOG-END-OFFSET  = final da fila, seria 30
	- LAG = a diferença do current pro log-end
	- CONSUMER-ID = dã                                   
	- HOST = dã           
	- CLIENT-ID = dã

```
GROUP           TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG             CONSUMER-ID                                       HOST            CLIENT-ID
x               teste           0          25              25              0               consumer-x-1-2cd8dac6-3e77-49bd-be94-4df5734c9068 /172.21.0.3     consumer-x-1
x               teste           1          31              31              0               consumer-x-1-2cd8dac6-3e77-49bd-be94-4df5734c9068 /172.21.0.3     consumer-x-1
x               teste           2          25              25              0               consumer-x-1-ec038717-2584-4c6f-b9e7-222a1f6d2021 /172.21.0.3     consumer-x-1
```

- 1 consumer tá lendo 2 partições

#control center
http://localhost:9021/