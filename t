saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java 2018-08-20T04:15:28,099Z [ATERT_104] INFO  s.f.p.e.ProcessingTaskExecutor  - Run complete, reset thread name from current REGISTRER-VIGSEL-DSF-OPPDATERT_104 to pool-1-thread-7
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java 2018-08-20T04:15:28,099Z [OTTATT_98] INFO  s.f.p.t.ProcessingTaskStateHandler  - Endrer state til DONE for task VIGSEL-MELDING-MOTTATT
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java 2018-08-20T04:15:28,100Z [OTTATT_98] INFO  s.f.p.t.ProcessingTaskStateHandler  -                                                                                                                                                                                                                    Task er ferdig [VIGSEL-MELDING-MOTTATT]1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java 2018-08-20T04:15:28,100Z [OTTATT_98] WARN  s.f.p.e.ProcessingTaskExecutor  - Retryable feil i task [VIGSEL-MELDING-MOTTATT], status = ske.folkeregister.processing.task.ProcessingTaskStatus@1e7923f[state=DONE,lastInput=,lastOutput=,numberOfElementsProcessed=0,skipStatus=ske.folkeregister.processing.task.SkipStatus@72355e8c[skipData=ske.folkeregister.processing.task.SkipData@4845e284[skip=[]],skipped=[]],subscriptionName=FREG-PROEVING-AV-EKTESKAPVIGSEL-MELDING-MOTTATT,processorMetaData=ske.folkeregister.processing.saksbehandling.SaksbehandlingProcessorMetaData@4ef96baf], retryCount = 523
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java ske.folkeregister.processing.execution.RetryableException: Feilkode: EVENT_PROCESSING_HENDELSENAV_2, melding: Uventet (ikke 200) response fra hendelsesnavet. Statuskode: [500], melding: [{"httpStatus":500,"alvorlighetsgrad":"ERROR","feilkode":"UNDEFINED","feilmelding":"Det oppstod en feil: Checked exception occurred while calling PDP.","tilleggsinformasjon":{"AURORA-transactionId":"c1003f87-bccf-427e-8377-1efd80c24ec6","AURORA-uri":"/hendelsesnav/api/subscription/info/FREG-PROEVING-AV-EKTESKAPVIGSEL-MELDING-MOTTATT","AURORA-system":"CN=folkeregistrering.saksflyt-vigsel","MDC-system":"CN=folkeregistrering.saksflyt-vigsel","MDC-ACCESS_CONTROL_ID":"CN=folkeregistrering.saksflyt-vigsel","MDC-Korrelasjonsid":"bbe04140-efc4-4a6d-a52a-a23b33612321","MDC-Klientid":"saksflyt-vigsel-proeving/feature_MFU_4720_Forbedret_brevskriving-SNAPSHOT","MDC-Meldingsid":"2d29936a-71cd-4c54-8408-83e845adf1d9","MDC-uri":"/hendelsesnav/api/subscription/info/FREG-PROEVING-AV-EKTESKAPVIGSEL-MELDING-MOTTATT","MDC-transactionId":"c1003f87-bccf-427e-8377-1efd80c24ec6"}}], kontekst:{}
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.infrastructure.hendelsesnavet.HttpHendelsesnavClient.httpStatusNotOkError(HttpHendelsesnavClient.java:127)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.infrastructure.hendelsesnavet.RequestBuilder.lambda$execute$0(RequestBuilder.java:103)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at org.apache.http.client.fluent.Response.handleResponse(Response.java:90)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.application.http.Request.sendOgHaandterResponse(Request.java:462)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.infrastructure.hendelsesnavet.RequestBuilder.execute(RequestBuilder.java:101)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.infrastructure.hendelsesnavet.HttpHendelsesnavClient.subscriptionsFor(HttpHendelsesnavClient.java:253)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.infrastructure.hendelsesnavet.PersistentFeedPageFetcherDefault.subscriptionInfo(PersistentFeedPageFetcherDefault.java:85)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.saksbehandling.SaksbehandlingProcessorSubscription.getStartPointerFor(SaksbehandlingProcessorSubscription.java:30)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.saksbehandling.SaksbehandlingProcessorSubscription.getStartPointerFor(SaksbehandlingProcessorSubscription.java:14)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.streaming.StreamingProcessingTask.sourceStreamForProcessor(StreamingProcessingTask.java:189)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.streaming.StreamingProcessingTask.resume(StreamingProcessingTask.java:183)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.streaming.StreamingProcessingTask.start(StreamingProcessingTask.java:74)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at ske.folkeregister.processing.execution.ProcessingTaskExecutor$RetryableProcessingTask.run(ProcessingTaskExecutor.java:189)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.util.concurrent.Executors$RunnableAdapter.call(Executors.java:511)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.util.concurrent.FutureTask.run(FutureTask.java:266)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.util.concurrent.ScheduledThreadPoolExecutor$ScheduledFutureTask.access$201(ScheduledThreadPoolExecutor.java:180)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.util.concurrent.ScheduledThreadPoolExecutor$ScheduledFutureTask.run(ScheduledThreadPoolExecutor.java:293)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1149)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:624)
saksflyt-vigsel-proeving-241-rxs7h saksflyt-vigsel-proeving-java                at java.lang.Thread.run(Thread.java:748)
