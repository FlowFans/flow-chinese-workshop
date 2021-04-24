//部署合约的Hash: 07ec92abd3496ecd8b2efd0854f3ba677fcf21dd530779af77bf6856f02903b7

//调用合约的Hash: a86569bf644a3b52efa4a374fd24946ba21a2483e15e0b643e79e4123c5d0df1


//合约代码:
 
access(all) contract HelloWorld2 {
  pub let greeting: String
  pub let question: String
  pub event HelloEvent(message: String)
  pub event AskQuesiton(message: String)

  init() {
    self.greeting = "Hello, World!"
    self.question = "How's your day?"
  }

  pub fun hello(message: String): String {
    emit HelloEvent(message: message)
    return self.greeting
  }

  pub fun askquestion(message: String): String {
    emit AskQuesiton(message: message)
    return self.question
  }
}
