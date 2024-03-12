using System;
using System.Collections.Generic;
using System.Diagnostics.Eventing.Reader;
using System.IO;
using System.Linq;
using System.Runtime.Remoting.Messaging;
using static System.Net.Mime.MediaTypeNames;


namespace log2accdb
{
    internal class Program
    {

        static Dictionary<string, string> resultsData = new Dictionary<string, string>();
        static StreamWriter writer;
        static StreamReader reader;
        static string fileOut = "out.txt";
        static string path = ".";
        static bool resultIPADDR;

        static void Main(string[] args)
        {
            if (args.Length < 2)
            {
                Console.WriteLine("Wrong usage!\n" + "log2db <search line> -s (result strings) or -a (result ip address)");
                return;
            }

            if (args[1] == "-s")
            {
                resultIPADDR = false;
            }
            else if (args[1] == "-a") 
            {
                resultIPADDR= true;
            }
            else
            {
                Console.WriteLine($"Key {args[1]} not define!\n");
                return;
            }

            string searchLine = args[0];
            writer = new StreamWriter(fileOut, false);



            //foreach (string arg in args)
            //    Console.WriteLine($"args: {arg} \n");

            if (File.Exists(path))
                ProcessFile(path, searchLine, false);

            else if (Directory.Exists(path))
                    ProcessDirectory(path, searchLine);

            else
                Console.WriteLine("{0} is not a valid file or directory.", path);

            


            if (!resultIPADDR)
                GetSessionID();

            writer.Close();

        }

        static void ProcessDirectory(string targetDirectory, string searchLine)
        {
            string[] fileEntries = Directory.GetFiles(targetDirectory);
            foreach (string fileName in fileEntries)
                ProcessFile(fileName, searchLine, false);
        }

        static void ProcessFile(string path, string searchLine, bool isWriteResults)
        {
            //long length = new System.IO.FileInfo(path).Length;
            //if (length > 0)
                ReadFromFile(path, searchLine, isWriteResults);

            //Console.WriteLine($"Processed file {path} size {length}");

        }


        static void ReadFromFile(string fileName, string searchLine, bool isWriteResults)
        {
            if (fileName.Contains("exe") || fileName.Contains(fileOut)) return;

            try
            {
                FileStream fs = new FileStream(fileName, FileMode.Open, FileAccess.Read, FileShare.ReadWrite);
                reader = new StreamReader(fs);

                while (!reader.EndOfStream)
                {
                    string s = reader.ReadLine();
                    if (!string.Equals(s[0], '#') && s.Contains(searchLine))
                    {
                        //Console.WriteLine($"read symbol {s[0]} from {fileName}");

                        if (!isWriteResults)
                        {
                            string[] fields = new string[7];
                            fields = s.Split(',');

                            if (resultIPADDR && string.Equals(fields[6], ">"))
                            {
                                string[] tmp = new string[1];
                                tmp = fields[5].Split(':');
                                string res = fields[0] + "," + tmp[0];

                                Console.WriteLine($"string to write: {res} ");
                                writer.WriteLine(res);

                            }
                            else if (!resultIPADDR)
                                resultsData.Add(fields[2], fileName);

                        }
                        else
                            writer.WriteLine(s);
                        

                    }

                }
                reader.Close();

            }
            catch (Exception e)
            {
                Console.WriteLine(e.ToString());
            }
     

        }


        static void GetSessionID()
        {
            
            foreach (var result in resultsData)
            {
                //System.Console.WriteLine($"{result.Key}, {result.Value}");
                writer.WriteLine("=============================================");
                ProcessFile(result.Value, result.Key, true);
            }

        }

    }
}

